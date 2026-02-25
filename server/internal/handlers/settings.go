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
func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Check admin role
	userID, _ := mw.GetUserID(r.Context())
	var role string
	err := h.DB.Pool.QueryRow(r.Context(), `SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to check user role")
		return
	}
	if role != "admin" {
		mw.Error(w, http.StatusForbidden, "admin access required")
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
