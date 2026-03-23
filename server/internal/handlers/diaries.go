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

type DiaryHandler struct {
	DB *database.DB
}

// GET /api/diaries
func (h *DiaryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, title, content, mood, weather, tags, is_private, created_at::text, updated_at::text
		FROM diaries WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query diaries")
		return
	}
	defer rows.Close()

	type diaryResp struct {
		ID        int64    `json:"id"`
		UserID    int64    `json:"user_id"`
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Mood      string   `json:"mood"`
		Weather   string   `json:"weather"`
		Tags      []string `json:"tags"`
		IsPrivate bool     `json:"is_private"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
	}

	var diaries []diaryResp
	for rows.Next() {
		var d diaryResp
		if err := rows.Scan(&d.ID, &d.UserID, &d.Title, &d.Content, &d.Mood,
			&d.Weather, &d.Tags, &d.IsPrivate, &d.CreatedAt, &d.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan diary")
			return
		}
		if d.Tags == nil {
			d.Tags = []string{}
		}
		diaries = append(diaries, d)
	}
	if diaries == nil {
		diaries = []diaryResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": diaries})
}

// POST /api/diaries
func (h *DiaryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Mood      string   `json:"mood"`
		Weather   string   `json:"weather"`
		Tags      []string `json:"tags"`
		IsPrivate bool     `json:"is_private"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Tags == nil {
		body.Tags = []string{}
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO diaries (user_id, title, content, mood, weather, tags, is_private)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, body.Title, body.Content, body.Mood, body.Weather, body.Tags, body.IsPrivate).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create diary")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "diary created"})
}

// GET /api/diaries/{id}
func (h *DiaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var d struct {
		ID        int64    `json:"id"`
		UserID    int64    `json:"user_id"`
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Mood      string   `json:"mood"`
		Weather   string   `json:"weather"`
		Tags      []string `json:"tags"`
		IsPrivate bool     `json:"is_private"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, user_id, title, content, mood, weather, tags, is_private, created_at::text, updated_at::text
		FROM diaries WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&d.ID, &d.UserID, &d.Title, &d.Content, &d.Mood,
		&d.Weather, &d.Tags, &d.IsPrivate, &d.CreatedAt, &d.UpdatedAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "diary not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get diary")
		return
	}
	if d.Tags == nil {
		d.Tags = []string{}
	}

	mw.JSON(w, http.StatusOK, d)
}

// PUT /api/diaries/{id}
func (h *DiaryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Mood      string   `json:"mood"`
		Weather   string   `json:"weather"`
		Tags      []string `json:"tags"`
		IsPrivate bool     `json:"is_private"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Tags == nil {
		body.Tags = []string{}
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE diaries SET title=$1, content=$2, mood=$3, weather=$4, tags=$5, is_private=$6, updated_at=NOW()
		WHERE id=$7 AND user_id=$8
	`, body.Title, body.Content, body.Mood, body.Weather, body.Tags, body.IsPrivate, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update diary")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "diary not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "diary updated"})
}

// DELETE /api/diaries/{id}
func (h *DiaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM diaries WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete diary")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "diary not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "diary deleted"})
}
