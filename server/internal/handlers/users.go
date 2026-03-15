package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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
		       is_banned, ban_reason, banned_at,
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

	var users []userResp
	for rows.Next() {
		var u userResp
		if err := rows.Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName,
			&u.AvatarURL, &u.Role, &u.IsBanned, &u.BanReason, &u.BannedAt,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
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
		       is_banned, ban_reason, banned_at,
		       created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName,
		&u.AvatarURL, &u.Role, &u.IsBanned, &u.BanReason, &u.BannedAt,
		&u.CreatedAt, &u.UpdatedAt)
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

	// Cannot demote yourself
	if id == callerID {
		mw.Error(w, http.StatusBadRequest, "cannot modify your own super_admin status")
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
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE users SET is_banned = true, ban_reason = $1, banned_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`, body.Reason, id)
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
