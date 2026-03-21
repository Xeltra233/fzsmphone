package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

func decodeCharacterExtra(raw []byte) map[string]interface{} {
	if len(raw) == 0 {
		return map[string]interface{}{}
	}

	var extra map[string]interface{}
	if err := json.Unmarshal(raw, &extra); err != nil || extra == nil {
		return map[string]interface{}{}
	}

	return extra
}

func estimateCharacterPayloadSize(name, avatarURL, description, personality, systemPrompt, greeting string, tags []string, extra map[string]interface{}) int64 {
	extraBytes, _ := json.Marshal(extra)
	tagsBytes, _ := json.Marshal(tags)
	return int64(len(name) + len(avatarURL) + len(description) + len(personality) + len(systemPrompt) + len(greeting) + len(tagsBytes) + len(extraBytes))
}

func (h *CharacterHandler) getCharacterStorageQuota(ctx context.Context, userID int64) (int64, error) {
	var quota int64
	err := h.DB.Pool.QueryRow(ctx, `SELECT COALESCE(character_storage_quota, 10485760) FROM users WHERE id = $1`, userID).Scan(&quota)
	return quota, err
}

func (h *CharacterHandler) getCharacterStorageUsed(ctx context.Context, userID int64) (int64, error) {
	var used int64
	err := h.DB.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(
			octet_length(name) +
			octet_length(avatar_url) +
			octet_length(description) +
			octet_length(personality) +
			octet_length(system_prompt) +
			octet_length(greeting) +
			octet_length(array_to_string(tags, ',')) +
			octet_length(extra::text)
		), 0)
		FROM characters WHERE user_id = $1
	`, userID).Scan(&used)
	return used, err
}

func (h *CharacterHandler) ensureCharacterQuota(ctx context.Context, userID int64, nextSize int64, excludeID int64) error {
	quota, err := h.getCharacterStorageQuota(ctx, userID)
	if err != nil {
		return err
	}

	var used int64
	if excludeID > 0 {
		err = h.DB.Pool.QueryRow(ctx, `
			SELECT COALESCE(SUM(
				octet_length(name) +
				octet_length(avatar_url) +
				octet_length(description) +
				octet_length(personality) +
				octet_length(system_prompt) +
				octet_length(greeting) +
				octet_length(array_to_string(tags, ',')) +
				octet_length(extra::text)
			), 0)
			FROM characters WHERE user_id = $1 AND id <> $2
		`, userID, excludeID).Scan(&used)
	} else {
		used, err = h.getCharacterStorageUsed(ctx, userID)
	}
	if err != nil {
		return err
	}

	if used+nextSize > quota {
		return fmt.Errorf("character storage quota exceeded")
	}
	return nil
}

// GET /api/characters
func (h *CharacterHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	search := r.URL.Query().Get("search")
	tag := r.URL.Query().Get("tag")

	query := `
		SELECT id, user_id, name, avatar_url, description, personality, system_prompt,
		       greeting, is_public, tags, extra, created_at, updated_at
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
		ID           int64                  `json:"id"`
		UserID       int64                  `json:"user_id"`
		Name         string                 `json:"name"`
		AvatarURL    string                 `json:"avatar_url"`
		Description  string                 `json:"description"`
		Personality  string                 `json:"personality"`
		SystemPrompt string                 `json:"system_prompt"`
		Greeting     string                 `json:"greeting"`
		IsPublic     bool                   `json:"is_public"`
		Tags         []string               `json:"tags"`
		Extra        map[string]interface{} `json:"extra"`
		CreatedAt    string                 `json:"created_at"`
		UpdatedAt    string                 `json:"updated_at"`
	}

	var characters []charResp
	for rows.Next() {
		var c charResp
		var extraRaw []byte
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.AvatarURL, &c.Description,
			&c.Personality, &c.SystemPrompt, &c.Greeting, &c.IsPublic, &c.Tags, &extraRaw,
			&c.CreatedAt, &c.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan character")
			return
		}
		if c.Tags == nil {
			c.Tags = []string{}
		}
		c.Extra = decodeCharacterExtra(extraRaw)
		characters = append(characters, c)
	}

	if characters == nil {
		characters = []charResp{}
	}

	quota, err := h.getCharacterStorageQuota(r.Context(), userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to load character quota")
		return
	}
	used, err := h.getCharacterStorageUsed(r.Context(), userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to load character storage usage")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"data": characters,
		"storage": map[string]int64{
			"quota_bytes": quota,
			"used_bytes":  used,
		},
	})
}

// POST /api/characters
func (h *CharacterHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Name         string                 `json:"name"`
		AvatarURL    string                 `json:"avatar_url"`
		Description  string                 `json:"description"`
		Personality  string                 `json:"personality"`
		SystemPrompt string                 `json:"system_prompt"`
		Greeting     string                 `json:"greeting"`
		IsPublic     bool                   `json:"is_public"`
		Tags         []string               `json:"tags"`
		Extra        map[string]interface{} `json:"extra"`
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
	if body.Extra == nil {
		body.Extra = map[string]interface{}{}
	}

	if err := h.ensureCharacterQuota(r.Context(), userID, estimateCharacterPayloadSize(body.Name, body.AvatarURL, body.Description, body.Personality, body.SystemPrompt, body.Greeting, body.Tags, body.Extra), 0); err != nil {
		if err.Error() == "character storage quota exceeded" {
			mw.Error(w, http.StatusBadRequest, "character storage quota exceeded")
			return
		}
		mw.Error(w, http.StatusInternalServerError, "failed to validate character storage quota")
		return
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO characters (user_id, name, avatar_url, description, personality, system_prompt, greeting, is_public, tags, extra)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, userID, body.Name, body.AvatarURL, body.Description, body.Personality,
		body.SystemPrompt, body.Greeting, body.IsPublic, body.Tags, body.Extra).Scan(&id)
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
		ID           int64                  `json:"id"`
		UserID       int64                  `json:"user_id"`
		Name         string                 `json:"name"`
		AvatarURL    string                 `json:"avatar_url"`
		Description  string                 `json:"description"`
		Personality  string                 `json:"personality"`
		SystemPrompt string                 `json:"system_prompt"`
		Greeting     string                 `json:"greeting"`
		IsPublic     bool                   `json:"is_public"`
		Tags         []string               `json:"tags"`
		Extra        map[string]interface{} `json:"extra"`
		CreatedAt    string                 `json:"created_at"`
		UpdatedAt    string                 `json:"updated_at"`
	}
	var extraRaw []byte

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, user_id, name, avatar_url, description, personality, system_prompt,
		       greeting, is_public, tags, extra, created_at, updated_at
		FROM characters WHERE id = $1 AND (user_id = $2 OR is_public = true)
	`, id, userID).Scan(&c.ID, &c.UserID, &c.Name, &c.AvatarURL, &c.Description,
		&c.Personality, &c.SystemPrompt, &c.Greeting, &c.IsPublic, &c.Tags, &extraRaw,
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
	c.Extra = decodeCharacterExtra(extraRaw)

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
		Name         string                 `json:"name"`
		AvatarURL    string                 `json:"avatar_url"`
		Description  string                 `json:"description"`
		Personality  string                 `json:"personality"`
		SystemPrompt string                 `json:"system_prompt"`
		Greeting     string                 `json:"greeting"`
		IsPublic     bool                   `json:"is_public"`
		Tags         []string               `json:"tags"`
		Extra        map[string]interface{} `json:"extra"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Tags == nil {
		body.Tags = []string{}
	}
	if body.Extra == nil {
		body.Extra = map[string]interface{}{}
	}

	if err := h.ensureCharacterQuota(r.Context(), userID, estimateCharacterPayloadSize(body.Name, body.AvatarURL, body.Description, body.Personality, body.SystemPrompt, body.Greeting, body.Tags, body.Extra), id); err != nil {
		if err.Error() == "character storage quota exceeded" {
			mw.Error(w, http.StatusBadRequest, "character storage quota exceeded")
			return
		}
		mw.Error(w, http.StatusInternalServerError, "failed to validate character storage quota")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE characters SET
			name = $1, avatar_url = $2, description = $3, personality = $4,
			system_prompt = $5, greeting = $6, is_public = $7, tags = $8, extra = $9, updated_at = NOW()
		WHERE id = $10 AND user_id = $11
	`, body.Name, body.AvatarURL, body.Description, body.Personality,
		body.SystemPrompt, body.Greeting, body.IsPublic, body.Tags, body.Extra, id, userID)
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
