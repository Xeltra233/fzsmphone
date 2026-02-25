package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fzsmphone/internal/config"
	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	DB  *database.DB
	Cfg *config.Config
}

type discordTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type discordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Discriminator string `json:"discriminator"`
}

// POST /api/auth/discord
func (h *AuthHandler) DiscordCallback(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		mw.Error(w, http.StatusBadRequest, "missing code")
		return
	}

	// Exchange code for token
	data := url.Values{
		"client_id":     {h.Cfg.DiscordClientID},
		"client_secret": {h.Cfg.DiscordSecret},
		"grant_type":    {"authorization_code"},
		"code":          {body.Code},
		"redirect_uri":  {h.Cfg.DiscordRedirect},
	}

	resp, err := http.Post(
		"https://discord.com/api/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		mw.Error(w, http.StatusBadGateway, "failed to contact Discord")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		mw.Error(w, http.StatusBadGateway, fmt.Sprintf("Discord token error: %s", string(body)))
		return
	}

	var tokenResp discordTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to parse Discord token")
		return
	}

	// Get user info
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := http.DefaultClient.Do(req)
	if err != nil {
		mw.Error(w, http.StatusBadGateway, "failed to get Discord user")
		return
	}
	defer userResp.Body.Close()

	var dUser discordUser
	if err := json.NewDecoder(userResp.Body).Decode(&dUser); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to parse Discord user")
		return
	}

	// Build avatar URL
	avatarURL := ""
	if dUser.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=256", dUser.ID, dUser.Avatar)
	}

	displayName := dUser.GlobalName
	if displayName == "" {
		displayName = dUser.Username
	}

	// Upsert user in database
	ctx := context.Background()
	var userID int64
	var role string

	err = h.DB.Pool.QueryRow(ctx, `
		INSERT INTO users (discord_id, username, display_name, avatar_url, email)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (discord_id) DO UPDATE SET
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			email = EXCLUDED.email,
			updated_at = NOW()
		RETURNING id, role
	`, dUser.ID, dUser.Username, displayName, avatarURL, dUser.Email).Scan(&userID, &role)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to save user")
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"discord_id": dUser.ID,
		"role":       role,
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenStr,
		"user": map[string]interface{}{
			"id":           userID,
			"discord_id":   dUser.ID,
			"username":     dUser.Username,
			"display_name": displayName,
			"avatar_url":   avatarURL,
			"role":         role,
		},
	})
}

// GET /api/auth/me
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		mw.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var user struct {
		ID          int64  `json:"id"`
		DiscordID   string `json:"discord_id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
		Role        string `json:"role"`
	}

	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, discord_id, username, display_name, avatar_url, role
		FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.DiscordID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.Role)
	if err != nil {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, user)
}
