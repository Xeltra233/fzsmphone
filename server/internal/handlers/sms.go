package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type SmsHandler struct {
	DB *database.DB
}

// GET /api/sms/threads
func (h *SmsHandler) ListThreads(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, recipient, number, character_id, avatar, last_content, last_at, created_at
		FROM sms_threads WHERE user_id = $1
		ORDER BY last_at DESC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query threads")
		return
	}
	defer rows.Close()

	type threadResp struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		Recipient   string `json:"recipient"`
		Number      string `json:"number"`
		CharacterID string `json:"character_id"`
		Avatar      string `json:"avatar"`
		LastContent string `json:"last_content"`
		LastAt      string `json:"last_at"`
		CreatedAt   string `json:"created_at"`
	}

	var threads []threadResp
	for rows.Next() {
		var t threadResp
		if err := rows.Scan(&t.ID, &t.UserID, &t.Recipient, &t.Number,
			&t.CharacterID, &t.Avatar, &t.LastContent, &t.LastAt, &t.CreatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan thread")
			return
		}
		threads = append(threads, t)
	}
	if threads == nil {
		threads = []threadResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": threads})
}

// POST /api/sms/threads
func (h *SmsHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Recipient   string `json:"recipient"`
		Number      string `json:"number"`
		CharacterID string `json:"character_id"`
		Avatar      string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO sms_threads (user_id, recipient, number, character_id, avatar)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, body.Recipient, body.Number, body.CharacterID, body.Avatar).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create thread")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "thread created"})
}

// DELETE /api/sms/threads/{id}
func (h *SmsHandler) DeleteThread(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM sms_threads WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete thread")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "thread not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "thread deleted"})
}

// GET /api/sms/threads/{id}/messages
func (h *SmsHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	threadID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verify thread ownership
	var ownerID int64
	err = h.DB.Pool.QueryRow(r.Context(), `SELECT user_id FROM sms_threads WHERE id = $1`, threadID).Scan(&ownerID)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "thread not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to verify thread")
		return
	}
	if ownerID != userID {
		mw.Error(w, http.StatusForbidden, "access denied")
		return
	}

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, thread_id, role, content, created_at
		FROM sms_messages WHERE thread_id = $1
		ORDER BY created_at ASC
	`, threadID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query messages")
		return
	}
	defer rows.Close()

	type msgResp struct {
		ID        int64  `json:"id"`
		ThreadID  int64  `json:"thread_id"`
		Role      string `json:"role"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}

	var messages []msgResp
	for rows.Next() {
		var m msgResp
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
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

// POST /api/sms/threads/{id}/messages
func (h *SmsHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	threadID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Verify thread ownership
	var ownerID int64
	err = h.DB.Pool.QueryRow(r.Context(), `SELECT user_id FROM sms_threads WHERE id = $1`, threadID).Scan(&ownerID)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "thread not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to verify thread")
		return
	}
	if ownerID != userID {
		mw.Error(w, http.StatusForbidden, "access denied")
		return
	}

	var body struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Role == "" {
		body.Role = "user"
	}

	var id int64
	err = h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO sms_messages (thread_id, role, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`, threadID, body.Role, body.Content).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create message")
		return
	}

	// Update thread last message
	_, _ = h.DB.Pool.Exec(r.Context(), `
		UPDATE sms_threads SET last_content = $1, last_at = NOW() WHERE id = $2
	`, body.Content, threadID)

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "message created"})
}
