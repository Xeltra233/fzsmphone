package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/jackc/pgx/v5"
)

type SettingsHandler struct {
	DB *database.DB
}

type APIProvider struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Provider  string      `json:"provider"`
	APIKey    string      `json:"api_key"`
	APIUrl    string      `json:"api_url"`
	CustomUrl string      `json:"custom_url,omitempty"`
	Model     string      `json:"model"`
	ModelList interface{} `json:"model_list,omitempty"`
	Enabled   bool        `json:"enabled"`
	ApiFormat string      `json:"api_format,omitempty"`
}

func normalizeProviderID(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return "provider"
	}
	value = strings.NewReplacer(" ", "-", "/", "-", "_", "-", ".", "-").Replace(value)
	return value
}

func inferProviderLabel(apiURL string) string {
	url := strings.ToLower(apiURL)
	switch {
	case strings.Contains(url, "openrouter"):
		return "OpenRouter"
	case strings.Contains(url, "deepseek"):
		return "DeepSeek"
	case strings.Contains(url, "generativelanguage.googleapis.com"):
		return "Gemini AI Studio"
	case strings.Contains(url, "openai"):
		return "OpenAI"
	case url == "":
		return "Custom"
	default:
		return "Custom"
	}
}

func normalizeProviderName(raw string, fallback string) string {
	raw = strings.TrimSpace(raw)
	if raw != "" {
		return raw
	}
	return fallback
}

func (h *SettingsHandler) getAppSettingProviders(r *http.Request, key string) []APIProvider {
	value, err := h.getAppSettingRaw(r, key)
	if err != nil {
		return nil
	}

	var providers []APIProvider
	if err := json.Unmarshal(value, &providers); err != nil {
		return nil
	}
	return providers
}

func (h *SettingsHandler) buildLegacyChatProvider(r *http.Request) APIProvider {
	apiURL := h.getAppSettingString(r, "apiUrl", "")
	return APIProvider{
		ID:        "chat-default",
		Name:      inferProviderLabel(apiURL),
		Provider:  normalizeProviderName(h.getAppSettingString(r, "provider", ""), "custom"),
		APIKey:    h.getAppSettingString(r, "apiKey", ""),
		APIUrl:    apiURL,
		CustomUrl: h.getAppSettingString(r, "custom_url", ""),
		Model:     h.getAppSettingString(r, "model", ""),
		ModelList: h.getAppSettingJSON(r, "model_list", []interface{}{}),
		Enabled:   true,
	}
}

func (h *SettingsHandler) buildLegacySocialProvider(r *http.Request) APIProvider {
	apiURL := h.getAppSettingString(r, "social_api_url", "")
	return APIProvider{
		ID:        "social-default",
		Name:      inferProviderLabel(apiURL),
		Provider:  normalizeProviderName(h.getAppSettingString(r, "social_provider", ""), "custom"),
		APIKey:    h.getAppSettingString(r, "social_api_key", ""),
		APIUrl:    apiURL,
		CustomUrl: h.getAppSettingString(r, "social_custom_url", ""),
		Model:     h.getAppSettingString(r, "social_model", ""),
		Enabled:   true,
	}
}

func (h *SettingsHandler) buildLegacyImageProviders(r *http.Request) []APIProvider {
	raw := h.getAppSettingString(r, "img_gen_config", "")
	if raw == "" {
		return nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return nil
	}

	providers := make([]APIProvider, 0, 3)
	for _, key := range []string{"openai", "novelai", "gemini"} {
		entry, ok := config[key].(map[string]interface{})
		if !ok {
			continue
		}
		apiURL, _ := entry["url"].(string)
		apiKey, _ := entry["key"].(string)
		model, _ := entry["model"].(string)
		if apiURL == "" && apiKey == "" && model == "" {
			continue
		}
		providers = append(providers, APIProvider{
			ID:        normalizeProviderID("image-" + key),
			Name:      strings.ToUpper(key),
			Provider:  key,
			APIKey:    apiKey,
			APIUrl:    apiURL,
			Model:     model,
			Enabled:   true,
			ApiFormat: key,
		})
	}
	return providers
}

func (h *SettingsHandler) getAppSettingRaw(r *http.Request, key string) ([]byte, error) {
	var value []byte
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT value FROM app_settings WHERE key = $1
	`, key).Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (h *SettingsHandler) getAppSettingString(r *http.Request, key, fallback string) string {
	value, err := h.getAppSettingRaw(r, key)
	if err != nil {
		return fallback
	}

	var parsed string
	if err := json.Unmarshal(value, &parsed); err == nil {
		return parsed
	}
	return string(value)
}

func (h *SettingsHandler) getAppSettingFloat(r *http.Request, key string, fallback float64) float64 {
	value, err := h.getAppSettingRaw(r, key)
	if err != nil {
		return fallback
	}

	var parsed float64
	if err := json.Unmarshal(value, &parsed); err == nil {
		return parsed
	}

	parsed, err = strconv.ParseFloat(strings.TrimSpace(string(value)), 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func (h *SettingsHandler) getAppSettingInt(r *http.Request, key string, fallback int) int {
	value, err := h.getAppSettingRaw(r, key)
	if err != nil {
		return fallback
	}

	var parsed int
	if err := json.Unmarshal(value, &parsed); err == nil {
		return parsed
	}

	parsed, err = strconv.Atoi(strings.TrimSpace(string(value)))
	if err != nil {
		return fallback
	}
	return parsed
}

func (h *SettingsHandler) getAppSettingJSON(r *http.Request, key string, fallback interface{}) interface{} {
	value, err := h.getAppSettingRaw(r, key)
	if err != nil {
		return fallback
	}

	var parsed interface{}
	if err := json.Unmarshal(value, &parsed); err != nil {
		return fallback
	}
	if parsed == nil {
		return fallback
	}
	return parsed
}

func (h *SettingsHandler) upsertAppSetting(r *http.Request, key string, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = h.DB.Pool.Exec(r.Context(), `
		INSERT INTO app_settings (key, value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()
	`, key, payload)
	return err
}

// GET /api/settings/public
// Returns safe public app settings for unauthenticated pages
func (h *SettingsHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	keys := []string{"app_name", "app_title", "announcement", "favicon"}
	rows, err := h.DB.Pool.Query(r.Context(), `
	SELECT key, value FROM app_settings WHERE key = ANY($1) ORDER BY key
	`, keys)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query public settings")
		return
	}
	defer rows.Close()

	settings := make(map[string]interface{})
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan public setting")
			return
		}

		var parsed interface{}
		if err := json.Unmarshal(value, &parsed); err == nil {
			settings[key] = parsed
		} else {
			settings[key] = string(value)
		}
	}

	mw.JSON(w, http.StatusOK, settings)
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

	settings := make(map[string]interface{})
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan setting")
			return
		}
		// Try to parse as JSON, if fails use as string
		var parsed interface{}
		if err := json.Unmarshal(value, &parsed); err == nil {
			settings[key] = parsed
		} else {
			settings[key] = string(value)
		}
	}

	mw.JSON(w, http.StatusOK, settings)
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
	availableModels := h.getAppSettingJSON(r, "model_list", []interface{}{})
	defaultModel := h.getAppSettingString(r, "model", "")

	// If super_admin, return global settings first, then user's personal settings
	if isSuperAdmin {
		var globalSettings struct {
			ID                      int64         `json:"id"`
			APIKey                  string        `json:"api_key"`
			APIUrl                  string        `json:"api_url"`
			CustomUrl               string        `json:"custom_url"`
			Model                   string        `json:"model"`
			ModelList               interface{}   `json:"model_list"`
			SocialAPIKey            string        `json:"social_api_key"`
			SocialAPIUrl            string        `json:"social_api_url"`
			SocialCustomUrl         string        `json:"social_custom_url"`
			SocialModel             string        `json:"social_model"`
			Temperature             float64       `json:"temperature"`
			MaxLength               int           `json:"max_length"`
			ContextSize             int           `json:"context_size"`
			Timeout                 int           `json:"timeout"`
			ChatProviders           []APIProvider `json:"chat_providers"`
			ChatDefaultProviderID   string        `json:"chat_default_provider_id"`
			SocialProviders         []APIProvider `json:"social_providers"`
			SocialDefaultProviderID string        `json:"social_default_provider_id"`
			ImageProviders          []APIProvider `json:"image_providers"`
			ImageDefaultProviderID  string        `json:"image_default_provider_id"`
		}
		err := h.DB.Pool.QueryRow(r.Context(), `
			SELECT id, COALESCE(api_key, ''), COALESCE(api_url, ''), COALESCE(model, '')
			FROM user_api_settings WHERE user_id = $1 AND is_global = true
		`, userID).Scan(&globalSettings.ID, &globalSettings.APIKey, &globalSettings.APIUrl, &globalSettings.Model)
		if err != nil && err != pgx.ErrNoRows {
			mw.Error(w, http.StatusInternalServerError, "failed to get global settings")
			return
		}

		if globalSettings.APIKey == "" {
			globalSettings.APIKey = h.getAppSettingString(r, "apiKey", "")
		}
		if globalSettings.APIUrl == "" {
			globalSettings.APIUrl = h.getAppSettingString(r, "apiUrl", "")
		}
		globalSettings.CustomUrl = h.getAppSettingString(r, "custom_url", "")
		if globalSettings.Model == "" {
			globalSettings.Model = h.getAppSettingString(r, "model", "")
		}

		globalSettings.ModelList = h.getAppSettingJSON(r, "model_list", []interface{}{})
		globalSettings.SocialAPIKey = h.getAppSettingString(r, "social_api_key", "")
		globalSettings.SocialAPIUrl = h.getAppSettingString(r, "social_api_url", "")
		globalSettings.SocialCustomUrl = h.getAppSettingString(r, "social_custom_url", "")
		globalSettings.SocialModel = h.getAppSettingString(r, "social_model", "")
		globalSettings.Temperature = h.getAppSettingFloat(r, "temperature", 0.9)
		globalSettings.MaxLength = h.getAppSettingInt(r, "max_length", 4000)
		globalSettings.ContextSize = h.getAppSettingInt(r, "context_size", 20)
		globalSettings.Timeout = h.getAppSettingInt(r, "timeout", 60)

		globalSettings.ChatProviders = h.getAppSettingProviders(r, "chat_providers")
		if len(globalSettings.ChatProviders) == 0 {
			globalSettings.ChatProviders = []APIProvider{h.buildLegacyChatProvider(r)}
		}
		globalSettings.ChatDefaultProviderID = h.getAppSettingString(r, "chat_default_provider_id", globalSettings.ChatProviders[0].ID)

		globalSettings.SocialProviders = h.getAppSettingProviders(r, "social_providers")
		if len(globalSettings.SocialProviders) == 0 {
			legacy := h.buildLegacySocialProvider(r)
			if legacy.APIUrl != "" || legacy.APIKey != "" || legacy.Model != "" {
				globalSettings.SocialProviders = []APIProvider{legacy}
			}
		}
		if len(globalSettings.SocialProviders) > 0 {
			globalSettings.SocialDefaultProviderID = h.getAppSettingString(r, "social_default_provider_id", globalSettings.SocialProviders[0].ID)
		}

		globalSettings.ImageProviders = h.getAppSettingProviders(r, "image_providers")
		if len(globalSettings.ImageProviders) == 0 {
			globalSettings.ImageProviders = h.buildLegacyImageProviders(r)
		}
		if len(globalSettings.ImageProviders) > 0 {
			globalSettings.ImageDefaultProviderID = h.getAppSettingString(r, "image_default_provider_id", globalSettings.ImageProviders[0].ID)
		}

		mw.JSON(w, http.StatusOK, map[string]interface{}{
			"is_super_admin":   true,
			"global":           globalSettings,
			"available_models": availableModels,
			"default_model":    defaultModel,
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
			"is_super_admin":   false,
			"personal":         nil,
			"available_models": availableModels,
			"default_model":    defaultModel,
		})
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get settings")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"is_super_admin":   false,
		"personal":         settings,
		"available_models": availableModels,
		"default_model":    defaultModel,
	})
}

// PUT /api/settings/api - Update user's API settings
func (h *SettingsHandler) UpdateUserApiSettings(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())

	var body struct {
		APIKey                  string          `json:"api_key"`
		APIUrl                  string          `json:"api_url"`
		CustomUrl               string          `json:"custom_url"`
		Model                   string          `json:"model"`
		ModelList               json.RawMessage `json:"model_list"`
		ChatProviders           json.RawMessage `json:"chat_providers"`
		ChatDefaultProviderID   string          `json:"chat_default_provider_id"`
		SocialAPIKey            string          `json:"social_api_key"`
		SocialAPIUrl            string          `json:"social_api_url"`
		SocialCustomUrl         string          `json:"social_custom_url"`
		SocialModel             string          `json:"social_model"`
		SocialProviders         json.RawMessage `json:"social_providers"`
		SocialDefaultProviderID string          `json:"social_default_provider_id"`
		ImageProviders          json.RawMessage `json:"image_providers"`
		ImageDefaultProviderID  string          `json:"image_default_provider_id"`
		Temperature             float64         `json:"temperature"`
		MaxLength               int             `json:"max_length"`
		ContextSize             int             `json:"context_size"`
		Timeout                 int             `json:"timeout"`
		IsGlobal                bool            `json:"is_global"` // Only super_admin can set global
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

	// Update first, then insert if missing. This avoids relying on a composite
	// unique index that may not exist in older deployments.
	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE user_api_settings
		SET api_key = $3, api_url = $4, model = $5, updated_at = NOW()
		WHERE user_id = $1 AND is_global = $2
	`, userID, body.IsGlobal, body.APIKey, body.APIUrl, body.Model)
	if err == nil && result.RowsAffected() == 0 {
		_, err = h.DB.Pool.Exec(r.Context(), `
			INSERT INTO user_api_settings (user_id, api_key, api_url, model, is_global, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`, userID, body.APIKey, body.APIUrl, body.Model, body.IsGlobal)
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update API settings")
		return
	}

	if body.IsGlobal {
		settingsToSync := map[string]interface{}{
			"apiKey":            body.APIKey,
			"apiUrl":            body.APIUrl,
			"custom_url":        body.CustomUrl,
			"model":             body.Model,
			"social_api_key":    body.SocialAPIKey,
			"social_api_url":    body.SocialAPIUrl,
			"social_custom_url": body.SocialCustomUrl,
			"social_model":      body.SocialModel,
			"temperature":       body.Temperature,
			"max_length":        body.MaxLength,
			"context_size":      body.ContextSize,
			"timeout":           body.Timeout,
		}

		for key, value := range settingsToSync {
			if err := h.upsertAppSetting(r, key, value); err != nil {
				mw.Error(w, http.StatusInternalServerError, "failed to update global setting: "+key)
				return
			}
		}

		modelListValue := interface{}([]interface{}{})
		if len(body.ModelList) > 0 {
			var parsed interface{}
			if err := json.Unmarshal(body.ModelList, &parsed); err != nil {
				mw.Error(w, http.StatusBadRequest, "invalid model_list")
				return
			}
			modelListValue = parsed
		}

		if err := h.upsertAppSetting(r, "model_list", modelListValue); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to update global setting: model_list")
			return
		}

		providerSettings := map[string]json.RawMessage{
			"chat_providers":   body.ChatProviders,
			"social_providers": body.SocialProviders,
			"image_providers":  body.ImageProviders,
		}
		for key, raw := range providerSettings {
			if len(raw) == 0 {
				continue
			}
			var parsed interface{}
			if err := json.Unmarshal(raw, &parsed); err != nil {
				mw.Error(w, http.StatusBadRequest, "invalid "+key)
				return
			}
			if err := h.upsertAppSetting(r, key, parsed); err != nil {
				mw.Error(w, http.StatusInternalServerError, "failed to update global setting: "+key)
				return
			}
		}

		for key, value := range map[string]string{
			"chat_default_provider_id":   body.ChatDefaultProviderID,
			"social_default_provider_id": body.SocialDefaultProviderID,
			"image_default_provider_id":  body.ImageDefaultProviderID,
		} {
			if err := h.upsertAppSetting(r, key, value); err != nil {
				mw.Error(w, http.StatusInternalServerError, "failed to update global setting: "+key)
				return
			}
		}
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "API settings updated"})
}

func (h *SettingsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID, isSuperAdmin := mw.GetUserID(r.Context())
	if userID == 0 {
		mw.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "only super admin can upload")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "failed to get file")
		return
	}
	defer file.Close()

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".ico": true, ".webp": true, ".svg": true}
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".png"
	}
	ext = "." + strings.TrimPrefix(strings.ToLower(ext), ".")
	if !allowedExts[ext] {
		mw.Error(w, http.StatusBadRequest, "invalid file type")
		return
	}

	newFilename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	uploadPath := filepath.Join(".", "public", newFilename)

	if !strings.HasPrefix(uploadPath, filepath.Join(".", "public")) {
		mw.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create directory")
		return
	}

	out, err := os.Create(uploadPath)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{
		"path": "/" + newFilename,
	})
}
