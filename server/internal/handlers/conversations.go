package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ConversationHandler struct {
	DB *database.DB
}

// GET /api/conversations
func (h *ConversationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT c.id, c.user_id, c.character_id, c.title, c.is_group, c.last_message,
		       c.last_at, c.created_at, c.updated_at,
		       ch.name AS char_name, ch.avatar_url AS char_avatar
		FROM conversations c
		LEFT JOIN characters ch ON ch.id = c.character_id
		WHERE c.user_id = $1
		ORDER BY c.last_at DESC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query conversations")
		return
	}
	defer rows.Close()

	type convResp struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		CharacterID int64  `json:"character_id"`
		Title       string `json:"title"`
		IsGroup     bool   `json:"is_group"`
		LastMessage string `json:"last_message"`
		LastAt      string `json:"last_at"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Character   *struct {
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
		} `json:"character,omitempty"`
	}

	var conversations []convResp
	for rows.Next() {
		var cv convResp
		var charName, charAvatar *string
		if err := rows.Scan(&cv.ID, &cv.UserID, &cv.CharacterID, &cv.Title, &cv.IsGroup,
			&cv.LastMessage, &cv.LastAt, &cv.CreatedAt, &cv.UpdatedAt,
			&charName, &charAvatar); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan conversation")
			return
		}
		if charName != nil {
			cv.Character = &struct {
				Name      string `json:"name"`
				AvatarURL string `json:"avatar_url"`
			}{Name: *charName, AvatarURL: ""}
			if charAvatar != nil {
				cv.Character.AvatarURL = *charAvatar
			}
		}
		conversations = append(conversations, cv)
	}

	if conversations == nil {
		conversations = []convResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": conversations})
}

// POST /api/conversations
func (h *ConversationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		CharacterID int64  `json:"character_id"`
		Title       string `json:"title"`
		IsGroup     bool   `json:"is_group"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO conversations (user_id, character_id, title, is_group, last_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, body.CharacterID, body.Title, body.IsGroup, time.Now()).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create conversation")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "conversation created",
	})
}

// GET /api/conversations/{id}
func (h *ConversationHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var cv struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		CharacterID int64  `json:"character_id"`
		Title       string `json:"title"`
		IsGroup     bool   `json:"is_group"`
		LastMessage string `json:"last_message"`
		LastAt      string `json:"last_at"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, user_id, character_id, title, is_group, last_message, last_at, created_at, updated_at
		FROM conversations WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&cv.ID, &cv.UserID, &cv.CharacterID, &cv.Title, &cv.IsGroup,
		&cv.LastMessage, &cv.LastAt, &cv.CreatedAt, &cv.UpdatedAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "conversation not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get conversation")
		return
	}

	mw.JSON(w, http.StatusOK, cv)
}

// DELETE /api/conversations/{id}
func (h *ConversationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Delete messages first, then conversation
	_, _ = h.DB.Pool.Exec(r.Context(), `DELETE FROM messages WHERE conversation_id = $1`, id)
	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM conversations WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete conversation")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "conversation not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "conversation deleted"})
}

// GET /api/conversations/{id}/messages
func (h *ConversationHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	convID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verify ownership
	var exists bool
	err = h.DB.Pool.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM conversations WHERE id = $1 AND user_id = $2)`,
		convID, userID).Scan(&exists)
	if err != nil || !exists {
		mw.Error(w, http.StatusNotFound, "conversation not found")
		return
	}

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, conversation_id, role, content, msg_type, created_at
		FROM messages WHERE conversation_id = $1
		ORDER BY created_at ASC
	`, convID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query messages")
		return
	}
	defer rows.Close()

	type msgResp struct {
		ID             int64  `json:"id"`
		ConversationID int64  `json:"conversation_id"`
		Role           string `json:"role"`
		Content        string `json:"content"`
		MsgType        string `json:"msg_type"`
		CreatedAt      string `json:"created_at"`
	}

	var messages []msgResp
	for rows.Next() {
		var m msgResp
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.MsgType, &m.CreatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan message")
			return
		}
		messages = append(messages, m)
	}

	if messages == nil {
		messages = []msgResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": messages})
}

// POST /api/conversations/{id}/messages
func (h *ConversationHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	convID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verify ownership
	var exists bool
	err = h.DB.Pool.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM conversations WHERE id = $1 AND user_id = $2)`,
		convID, userID).Scan(&exists)
	if err != nil || !exists {
		mw.Error(w, http.StatusNotFound, "conversation not found")
		return
	}

	var body struct {
		Role    string `json:"role"`
		Content string `json:"content"`
		MsgType string `json:"msg_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Role == "" {
		body.Role = "user"
	}
	if body.MsgType == "" {
		body.MsgType = "text"
	}

	var id int64
	err = h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO messages (conversation_id, role, content, msg_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, convID, body.Role, body.Content, body.MsgType).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create message")
		return
	}

	// Update conversation last_message
	truncated := body.Content
	if len(truncated) > 100 {
		truncated = truncated[:100] + "..."
	}
	_, _ = h.DB.Pool.Exec(r.Context(), `
		UPDATE conversations SET last_message = $1, last_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`, truncated, convID)

	mw.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "message created",
	})
}
