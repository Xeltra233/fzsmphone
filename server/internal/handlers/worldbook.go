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

type WorldBookHandler struct {
	DB *database.DB
}

// GET /api/worldbook
func (h *WorldBookHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, key, content, keywords, is_enabled, priority, created_at, updated_at
		FROM worldbook_entries WHERE user_id = $1
		ORDER BY priority DESC, created_at ASC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query worldbook entries")
		return
	}
	defer rows.Close()

	type entryResp struct {
		ID        int64    `json:"id"`
		UserID    int64    `json:"user_id"`
		Key       string   `json:"key"`
		Content   string   `json:"content"`
		Keywords  []string `json:"keywords"`
		IsEnabled bool     `json:"is_enabled"`
		Priority  int      `json:"priority"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
	}

	var entries []entryResp
	for rows.Next() {
		var e entryResp
		if err := rows.Scan(&e.ID, &e.UserID, &e.Key, &e.Content,
			&e.Keywords, &e.IsEnabled, &e.Priority, &e.CreatedAt, &e.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan worldbook entry")
			return
		}
		if e.Keywords == nil {
			e.Keywords = []string{}
		}
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []entryResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": entries})
}

// POST /api/worldbook
func (h *WorldBookHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Key       string   `json:"key"`
		Content   string   `json:"content"`
		Keywords  []string `json:"keywords"`
		IsEnabled bool     `json:"is_enabled"`
		Priority  int      `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Key == "" {
		mw.Error(w, http.StatusBadRequest, "key is required")
		return
	}
	if body.Keywords == nil {
		body.Keywords = []string{}
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO worldbook_entries (user_id, key, content, keywords, is_enabled, priority)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, userID, body.Key, body.Content, body.Keywords, body.IsEnabled, body.Priority).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create worldbook entry")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "worldbook entry created"})
}

// GET /api/worldbook/{id}
func (h *WorldBookHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var e struct {
		ID        int64    `json:"id"`
		UserID    int64    `json:"user_id"`
		Key       string   `json:"key"`
		Content   string   `json:"content"`
		Keywords  []string `json:"keywords"`
		IsEnabled bool     `json:"is_enabled"`
		Priority  int      `json:"priority"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, user_id, key, content, keywords, is_enabled, priority, created_at, updated_at
		FROM worldbook_entries WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&e.ID, &e.UserID, &e.Key, &e.Content,
		&e.Keywords, &e.IsEnabled, &e.Priority, &e.CreatedAt, &e.UpdatedAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "worldbook entry not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get worldbook entry")
		return
	}
	if e.Keywords == nil {
		e.Keywords = []string{}
	}

	mw.JSON(w, http.StatusOK, e)
}

// PUT /api/worldbook/{id}
func (h *WorldBookHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Key       string   `json:"key"`
		Content   string   `json:"content"`
		Keywords  []string `json:"keywords"`
		IsEnabled bool     `json:"is_enabled"`
		Priority  int      `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Keywords == nil {
		body.Keywords = []string{}
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE worldbook_entries SET key=$1, content=$2, keywords=$3, is_enabled=$4, priority=$5, updated_at=NOW()
		WHERE id=$6 AND user_id=$7
	`, body.Key, body.Content, body.Keywords, body.IsEnabled, body.Priority, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update worldbook entry")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "worldbook entry not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "worldbook entry updated"})
}

// DELETE /api/worldbook/{id}
func (h *WorldBookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM worldbook_entries WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete worldbook entry")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "worldbook entry not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "worldbook entry deleted"})
}
