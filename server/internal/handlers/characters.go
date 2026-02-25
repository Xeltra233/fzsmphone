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

type CharacterHandler struct {
	DB *database.DB
}

// GET /api/characters
func (h *CharacterHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	search := r.URL.Query().Get("search")
	tag := r.URL.Query().Get("tag")

	query := `
		SELECT id, user_id, name, avatar_url, description, personality, system_prompt,
		       greeting, is_public, tags, created_at, updated_at
		FROM characters
		WHERE (user_id = $1 OR is_public = true)
	`
	args := []interface{}{userID}
	argIdx := 2

	if search != "" {
		query += ` AND (name ILIKE $` + strconv.Itoa(argIdx) + ` OR description ILIKE $` + strconv.Itoa(argIdx) + `)`
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if tag != "" {
		query += ` AND $` + strconv.Itoa(argIdx) + ` = ANY(tags)`
		args = append(args, tag)
		argIdx++
	}

	query += ` ORDER BY updated_at DESC`

	rows, err := h.DB.Pool.Query(r.Context(), query, args...)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query characters")
		return
	}
	defer rows.Close()

	type charResp struct {
		ID           int64    `json:"id"`
		UserID       int64    `json:"user_id"`
		Name         string   `json:"name"`
		AvatarURL    string   `json:"avatar_url"`
		Description  string   `json:"description"`
		Personality  string   `json:"personality"`
		SystemPrompt string   `json:"system_prompt"`
		Greeting     string   `json:"greeting"`
		IsPublic     bool     `json:"is_public"`
		Tags         []string `json:"tags"`
		CreatedAt    string   `json:"created_at"`
		UpdatedAt    string   `json:"updated_at"`
	}

	var characters []charResp
	for rows.Next() {
		var c charResp
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.AvatarURL, &c.Description,
			&c.Personality, &c.SystemPrompt, &c.Greeting, &c.IsPublic, &c.Tags,
			&c.CreatedAt, &c.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan character")
			return
		}
		if c.Tags == nil {
			c.Tags = []string{}
		}
		characters = append(characters, c)
	}

	if characters == nil {
		characters = []charResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": characters})
}

// POST /api/characters
func (h *CharacterHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Name         string   `json:"name"`
		AvatarURL    string   `json:"avatar_url"`
		Description  string   `json:"description"`
		Personality  string   `json:"personality"`
		SystemPrompt string   `json:"system_prompt"`
		Greeting     string   `json:"greeting"`
		IsPublic     bool     `json:"is_public"`
		Tags         []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		mw.Error(w, http.StatusBadRequest, "name is required")
		return
	}
	if body.Tags == nil {
		body.Tags = []string{}
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO characters (user_id, name, avatar_url, description, personality, system_prompt, greeting, is_public, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, userID, body.Name, body.AvatarURL, body.Description, body.Personality,
		body.SystemPrompt, body.Greeting, body.IsPublic, body.Tags).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create character")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "character created",
	})
}

// GET /api/characters/{id}
func (h *CharacterHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var c struct {
		ID           int64    `json:"id"`
		UserID       int64    `json:"user_id"`
		Name         string   `json:"name"`
		AvatarURL    string   `json:"avatar_url"`
		Description  string   `json:"description"`
		Personality  string   `json:"personality"`
		SystemPrompt string   `json:"system_prompt"`
		Greeting     string   `json:"greeting"`
		IsPublic     bool     `json:"is_public"`
		Tags         []string `json:"tags"`
		CreatedAt    string   `json:"created_at"`
		UpdatedAt    string   `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, user_id, name, avatar_url, description, personality, system_prompt,
		       greeting, is_public, tags, created_at, updated_at
		FROM characters WHERE id = $1 AND (user_id = $2 OR is_public = true)
	`, id, userID).Scan(&c.ID, &c.UserID, &c.Name, &c.AvatarURL, &c.Description,
		&c.Personality, &c.SystemPrompt, &c.Greeting, &c.IsPublic, &c.Tags,
		&c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "character not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get character")
		return
	}
	if c.Tags == nil {
		c.Tags = []string{}
	}

	mw.JSON(w, http.StatusOK, c)
}

// PUT /api/characters/{id}
func (h *CharacterHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Name         string   `json:"name"`
		AvatarURL    string   `json:"avatar_url"`
		Description  string   `json:"description"`
		Personality  string   `json:"personality"`
		SystemPrompt string   `json:"system_prompt"`
		Greeting     string   `json:"greeting"`
		IsPublic     bool     `json:"is_public"`
		Tags         []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Tags == nil {
		body.Tags = []string{}
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE characters SET
			name = $1, avatar_url = $2, description = $3, personality = $4,
			system_prompt = $5, greeting = $6, is_public = $7, tags = $8, updated_at = NOW()
		WHERE id = $9 AND user_id = $10
	`, body.Name, body.AvatarURL, body.Description, body.Personality,
		body.SystemPrompt, body.Greeting, body.IsPublic, body.Tags, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update character")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "character not found or not owned by you")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "character updated"})
}

// DELETE /api/characters/{id}
func (h *CharacterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM characters WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete character")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "character not found or not owned by you")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "character deleted"})
}
