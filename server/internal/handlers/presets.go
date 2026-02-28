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

type PresetHandler struct {
	DB *database.DB
}

// GET /api/presets
func (h *PresetHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, name, emoji, category, description, content, prefill,
		       enable_prefill, is_builtin, gradient, created_at, updated_at
		FROM presets WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query presets")
		return
	}
	defer rows.Close()

	type presetResp struct {
		ID            int64  `json:"id"`
		UserID        int64  `json:"user_id"`
		Name          string `json:"name"`
		Emoji         string `json:"emoji"`
		Category      string `json:"category"`
		Description   string `json:"description"`
		Content       string `json:"content"`
		Prefill       string `json:"prefill"`
		EnablePrefill bool   `json:"enable_prefill"`
		IsBuiltin     bool   `json:"is_builtin"`
		Gradient      string `json:"gradient"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	}

	var presets []presetResp
	for rows.Next() {
		var p presetResp
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Emoji, &p.Category,
			&p.Description, &p.Content, &p.Prefill, &p.EnablePrefill,
			&p.IsBuiltin, &p.Gradient, &p.CreatedAt, &p.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan preset")
			return
		}
		presets = append(presets, p)
	}
	if presets == nil {
		presets = []presetResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": presets})
}

// POST /api/presets
func (h *PresetHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Name          string `json:"name"`
		Emoji         string `json:"emoji"`
		Category      string `json:"category"`
		Description   string `json:"description"`
		Content       string `json:"content"`
		Prefill       string `json:"prefill"`
		EnablePrefill bool   `json:"enable_prefill"`
		IsBuiltin     bool   `json:"is_builtin"`
		Gradient      string `json:"gradient"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Emoji == "" {
		body.Emoji = "◈"
	}
	if body.Gradient == "" {
		body.Gradient = "linear-gradient(135deg, #667eea, #764ba2)"
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO presets (user_id, name, emoji, category, description, content, prefill, enable_prefill, is_builtin, gradient)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, userID, body.Name, body.Emoji, body.Category, body.Description,
		body.Content, body.Prefill, body.EnablePrefill, body.IsBuiltin, body.Gradient).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create preset")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "preset created"})
}

// GET /api/presets/{id}
func (h *PresetHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p struct {
		ID            int64  `json:"id"`
		UserID        int64  `json:"user_id"`
		Name          string `json:"name"`
		Emoji         string `json:"emoji"`
		Category      string `json:"category"`
		Description   string `json:"description"`
		Content       string `json:"content"`
		Prefill       string `json:"prefill"`
		EnablePrefill bool   `json:"enable_prefill"`
		IsBuiltin     bool   `json:"is_builtin"`
		Gradient      string `json:"gradient"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, user_id, name, emoji, category, description, content, prefill,
		       enable_prefill, is_builtin, gradient, created_at, updated_at
		FROM presets WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&p.ID, &p.UserID, &p.Name, &p.Emoji, &p.Category,
		&p.Description, &p.Content, &p.Prefill, &p.EnablePrefill,
		&p.IsBuiltin, &p.Gradient, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "preset not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get preset")
		return
	}

	mw.JSON(w, http.StatusOK, p)
}

// PUT /api/presets/{id}
func (h *PresetHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Name          string `json:"name"`
		Emoji         string `json:"emoji"`
		Category      string `json:"category"`
		Description   string `json:"description"`
		Content       string `json:"content"`
		Prefill       string `json:"prefill"`
		EnablePrefill bool   `json:"enable_prefill"`
		Gradient      string `json:"gradient"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE presets SET name=$1, emoji=$2, category=$3, description=$4,
		       content=$5, prefill=$6, enable_prefill=$7, gradient=$8, updated_at=NOW()
		WHERE id=$9 AND user_id=$10
	`, body.Name, body.Emoji, body.Category, body.Description,
		body.Content, body.Prefill, body.EnablePrefill, body.Gradient, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update preset")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "preset not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "preset updated"})
}

// DELETE /api/presets/{id}
func (h *PresetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM presets WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete preset")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "preset not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "preset deleted"})
}
