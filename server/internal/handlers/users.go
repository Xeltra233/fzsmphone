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

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	DB *database.DB
}

// GET /api/users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(), `
	SELECT id, discord_id, username, display_name, COALESCE(avatar_url, ''), role,
	is_banned, COALESCE(ban_reason, ''), banned_at,
	COALESCE(credits, 0), COALESCE(total_tokens, 0), COALESCE(signin_streak, 0), COALESCE(invite_code, ''),
	created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	`)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query users")
		return
	}
	defer rows.Close()

	type userResp struct {
		ID           int64      `json:"id"`
		DiscordID    string     `json:"discord_id"`
		Username     string     `json:"username"`
		DisplayName  string     `json:"display_name"`
		AvatarURL    string     `json:"avatar_url"`
		Role         string     `json:"role"`
		IsBanned     bool       `json:"is_banned"`
		BanReason    string     `json:"ban_reason"`
		BannedAt     *time.Time `json:"banned_at"`
		Credits      int        `json:"credits"`
		TotalTokens  int        `json:"total_tokens"`
		SigninStreak int        `json:"signin_streak"`
		InviteCode   string     `json:"invite_code"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
	}

	var users []userResp
	for rows.Next() {
		var u userResp
		if err := rows.Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Role, &u.IsBanned, &u.BanReason, &u.BannedAt, &u.Credits, &u.TotalTokens, &u.SigninStreak, &u.InviteCode, &u.CreatedAt, &u.UpdatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan user")
			return
		}
		users = append(users, u)
	}
	if users == nil {
		users = []userResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": users})
}

// GET /api/users/{id}
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var u struct {
		ID          int64      `json:"id"`
		DiscordID   string     `json:"discord_id"`
		Username    string     `json:"username"`
		DisplayName string     `json:"display_name"`
		AvatarURL   string     `json:"avatar_url"`
		Role        string     `json:"role"`
		IsBanned    bool       `json:"is_banned"`
		BanReason   string     `json:"ban_reason"`
		BannedAt    *time.Time `json:"banned_at"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
	SELECT id, discord_id, username, display_name, COALESCE(avatar_url, ''), role,
	is_banned, COALESCE(ban_reason, ''), banned_at,
	created_at, updated_at WHERE id = $1
	`, id).Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Role, &u.IsBanned, &u.BanReason, &u.BannedAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, u)
}

// PATCH /api/users/{id}/role
func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Only super_admin can update roles
	_, _ = mw.GetUserID(r.Context())
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "super_admin access required")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Role != "admin" && body.Role != "user" && body.Role != "moderator" && body.Role != "super_admin" {
		mw.Error(w, http.StatusBadRequest, "invalid role")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
	UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2
	`, body.Role, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update role")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

// POST /api/users/{id}/set-super-admin
func (h *UserHandler) SetSuperAdmin(w http.ResponseWriter, r *http.Request) {
	// Only super_admin can set other super_admin
	callerID, _ := mw.GetUserID(r.Context())
	callerIsSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !callerIsSuperAdmin {
		mw.Error(w, http.StatusForbidden, "super_admin access required")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Cannot modify yourself
	if id == callerID {
		mw.Error(w, http.StatusBadRequest, "cannot modify your own super_admin status")
		return
	}

	// Get target user info to check privilege level
	var targetIsAdmin, targetIsSuperAdmin bool
	err = h.DB.Pool.QueryRow(r.Context(), `SELECT is_admin, is_super_admin FROM users WHERE id = $1`, id).Scan(&targetIsAdmin, &targetIsSuperAdmin)
	if err != nil {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	// Cannot modify users with same or higher privilege level
	if targetIsSuperAdmin || targetIsAdmin {
		mw.Error(w, http.StatusForbidden, "cannot modify users with same or higher privilege level")
		return
	}

	var body struct {
		IsSuperAdmin bool `json:"is_super_admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
	UPDATE users SET is_super_admin = $1, updated_at = NOW() WHERE id = $2
	`, body.IsSuperAdmin, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update super_admin status")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "super_admin status updated"})
}

// POST /api/users/{id}/ban
func (h *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
	// Only admins can ban users
	callerID, _ := mw.GetUserID(r.Context())
	var callerRole string
	err := h.DB.Pool.QueryRow(r.Context(), `SELECT role FROM users WHERE id = $1`, callerID).Scan(&callerRole)
	if err != nil || callerRole != "admin" {
		mw.Error(w, http.StatusForbidden, "admin access required")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Cannot ban yourself
	if id == callerID {
		mw.Error(w, http.StatusBadRequest, "cannot ban yourself")
		return
	}

	// Cannot ban another admin
	var targetRole string
	err = h.DB.Pool.QueryRow(r.Context(), `SELECT role FROM users WHERE id = $1`, id).Scan(&targetRole)
	if err != nil {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}
	if targetRole == "admin" {
		mw.Error(w, http.StatusForbidden, "cannot ban an admin")
		return
	}

	var body struct {
		Reason        string `json:"reason"`
		DurationHours int    `json:"duration_hours"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	var bannedUntil *time.Time
	if body.DurationHours > 0 {
		until := time.Now().Add(time.Duration(body.DurationHours) * time.Hour)
		bannedUntil = &until
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE users SET is_banned = true, ban_reason = $1, banned_at = NOW(), banned_until = $2, updated_at = NOW()
		WHERE id = $3
	`, body.Reason, bannedUntil, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to ban user")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "user banned"})
}

// POST /api/users/{id}/unban
func (h *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
	// Only admins can unban users
	callerID, _ := mw.GetUserID(r.Context())
	var callerRole string
	err := h.DB.Pool.QueryRow(r.Context(), `SELECT role FROM users WHERE id = $1`, callerID).Scan(&callerRole)
	if err != nil || callerRole != "admin" {
		mw.Error(w, http.StatusForbidden, "admin access required")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
	UPDATE users SET is_banned = false, ban_reason = '', banned_at = NULL, updated_at = NOW()
	WHERE id = $1
	`, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to unban user")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "user unbanned"})
}

// PATCH /api/users/{id}/credits
func (h *UserHandler) UpdateCredits(w http.ResponseWriter, r *http.Request) {
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "需要超级管理员权限")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Credits int `json:"credits"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
	UPDATE users SET credits = $1, updated_at = NOW() WHERE id = $2
	`, body.Credits, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update credits")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "credits updated"})
}

// PATCH /api/users/{id}/avatar - Upload user avatar
func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	callerID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if callerID != id {
		mw.Error(w, http.StatusForbidden, "can only update your own avatar")
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "failed to get file")
		return
	}
	defer file.Close()

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".png"
	}
	ext = "." + strings.TrimPrefix(strings.ToLower(ext), ".")
	if !allowedExts[ext] {
		mw.Error(w, http.StatusBadRequest, "invalid file type")
		return
	}

	newFilename := fmt.Sprintf("avatar_%d_%d%s", id, time.Now().Unix(), ext)
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

	avatarPath := "/" + newFilename
	_, err = h.DB.Pool.Exec(r.Context(), `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`, avatarPath, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update avatar")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"avatar_url": avatarPath, "message": "avatar updated"})
}

// PATCH /api/users/{id}/profile - Update user profile (display_name, avatar_url)
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	callerID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if callerID != id {
		mw.Error(w, http.StatusForbidden, "can only update your own profile")
		return
	}

	var body struct {
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.DisplayName != "" && len(body.DisplayName) > 50 {
		mw.Error(w, http.StatusBadRequest, "display name too long (max 50 chars)")
		return
	}
	if body.AvatarURL != "" && len(body.AvatarURL) > 500 {
		mw.Error(w, http.StatusBadRequest, "avatar URL too long (max 500 chars)")
		return
	}
	_, err = h.DB.Pool.Exec(r.Context(), `UPDATE users SET display_name = COALESCE(NULLIF($1, ''), display_name), avatar = COALESCE(NULLIF($2, ''), avatar), updated_at = NOW() WHERE id = $3`, body.DisplayName, body.AvatarURL, id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "profile updated", "avatar_url": body.AvatarURL})
}
