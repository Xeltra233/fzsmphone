package handlers

import (
	"encoding/json"
	"net/http"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/jackc/pgx/v5"
)

type SettingsHandler struct {
	DB *database.DB
}

// GET /api/settings
// Returns all app settings as a key-value map
func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(), `
	SELECT key, value FROM app_settings ORDER BY key
	`)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query settings")
		return
	}
	defer rows.Close()

	settings := make(map[string]json.RawMessage)
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan setting")
			return
		}
		settings[key] = json.RawMessage(value)
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": settings})
}

// PUT /api/settings
// Accepts a JSON object of key-value pairs and upserts them
// Only super_admin can update settings
func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	_, _ = mw.GetUserID(r.Context())
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "super_admin access required")
		return
	}

	var body map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for key, value := range body {
		_, err := h.DB.Pool.Exec(r.Context(), `
		INSERT INTO app_settings (key, value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()
		`, key, []byte(value))
		if err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to update setting: "+key)
			return
		}
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "settings updated"})
}

// GET /api/settings/{key}
// Returns a single setting by key
func (h *SettingsHandler) GetByKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		mw.Error(w, http.StatusBadRequest, "key parameter is required")
		return
	}

	var value []byte
	err := h.DB.Pool.QueryRow(r.Context(), `
	SELECT value FROM app_settings WHERE key = $1
	`, key).Scan(&value)
	if err == pgx.ErrNoRows {
		mw.JSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": nil,
		})
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get setting")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"key":   key,
		"value": json.RawMessage(value),
	})
}

// === User API Settings ===

// GET /api/settings/api - Get user's API settings (or global if super_admin)
func (h *SettingsHandler) GetUserApiSettings(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())

	// If super_admin, return global settings first, then user's personal settings
	if isSuperAdmin {
		var globalSettings struct {
			ID     int64  `json:"id"`
			APIKey string `json:"api_key"`
			APIUrl string `json:"api_url"`
			Model  string `json:"model"`
		}
		err := h.DB.Pool.QueryRow(r.Context(), `
			SELECT id, COALESCE(api_key, ''), COALESCE(api_url, ''), COALESCE(model, '')
			FROM user_api_settings WHERE user_id = $1 AND is_global = true
		`, userID).Scan(&globalSettings.ID, &globalSettings.APIKey, &globalSettings.APIUrl, &globalSettings.Model)
		if err != nil && err != pgx.ErrNoRows {
			mw.Error(w, http.StatusInternalServerError, "failed to get global settings")
			return
		}
		mw.JSON(w, http.StatusOK, map[string]interface{}{
			"is_super_admin": true,
			"global":         globalSettings,
		})
		return
	}

	// Regular users: only get their personal settings
	var settings struct {
		ID     int64  `json:"id"`
		APIKey string `json:"api_key"`
		APIUrl string `json:"api_url"`
		Model  string `json:"model"`
	}
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, COALESCE(api_key, ''), COALESCE(api_url, ''), COALESCE(model, '')
		FROM user_api_settings WHERE user_id = $1 AND is_global = false
	`, userID).Scan(&settings.ID, &settings.APIKey, &settings.APIUrl, &settings.Model)
	if err == pgx.ErrNoRows {
		mw.JSON(w, http.StatusOK, map[string]interface{}{
			"is_super_admin": false,
			"personal":       nil,
		})
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get settings")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"is_super_admin": false,
		"personal":       settings,
	})
}

// PUT /api/settings/api - Update user's API settings
func (h *SettingsHandler) UpdateUserApiSettings(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())

	var body struct {
		APIKey   string `json:"api_key"`
		APIUrl   string `json:"api_url"`
		Model    string `json:"model"`
		IsGlobal bool   `json:"is_global"` // Only super_admin can set global
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Check permission for global settings
	if body.IsGlobal && !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "only super_admin can set global API settings")
		return
	}

	// Upsert user API settings
	_, err := h.DB.Pool.Exec(r.Context(), `
		INSERT INTO user_api_settings (user_id, api_key, api_url, model, is_global, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (user_id, is_global) 
		WHERE is_global = $5
		DO UPDATE SET api_key = $2, api_url = $3, model = $4, updated_at = NOW()
	`, userID, body.APIKey, body.APIUrl, body.Model, body.IsGlobal)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update API settings")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "API settings updated"})
}
