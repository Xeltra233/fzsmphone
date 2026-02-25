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

type PersonaHandler struct {
	DB *database.DB
}

// GET /api/personas
func (h *PersonaHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, name, description, avatar_url, is_default, created_at, updated_at
		FROM user_personas WHERE user_id = $1
		ORDER BY is_default DESC, created_at ASC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query personas")
		return
	}
	defer rows.Close()

	type personaResp struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
		IsDefault   bool   `json:"is_default"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	var personas []personaResp
	for rows.Next() {
		var p personaResp
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description,
			&p.AvatarURL, &p.IsDefault, &p.CreatedAt, &p.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan persona")
			return
		}
		personas = append(personas, p)
	}
	if personas == nil {
		personas = []personaResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": personas})
}

// POST /api/personas
func (h *PersonaHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
		IsDefault   bool   `json:"is_default"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		mw.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	// If setting as default, clear other defaults first
	if body.IsDefault {
		_, _ = h.DB.Pool.Exec(r.Context(), `
			UPDATE user_personas SET is_default = false WHERE user_id = $1
		`, userID)
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO user_personas (user_id, name, description, avatar_url, is_default)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, body.Name, body.Description, body.AvatarURL, body.IsDefault).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create persona")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "persona created"})
}

// PUT /api/personas/{id}
func (h *PersonaHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
		IsDefault   bool   `json:"is_default"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// If setting as default, clear other defaults first
	if body.IsDefault {
		_, _ = h.DB.Pool.Exec(r.Context(), `
			UPDATE user_personas SET is_default = false WHERE user_id = $1 AND id != $2
		`, userID, id)
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE user_personas SET name=$1, description=$2, avatar_url=$3, is_default=$4, updated_at=NOW()
		WHERE id=$5 AND user_id=$6
	`, body.Name, body.Description, body.AvatarURL, body.IsDefault, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update persona")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "persona not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "persona updated"})
}

// DELETE /api/personas/{id}
func (h *PersonaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Don't allow deleting the default persona if it's the only one
	var isDefault bool
	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT is_default FROM user_personas WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&isDefault)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "persona not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to check persona")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM user_personas WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete persona")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "persona not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "persona deleted"})
}
