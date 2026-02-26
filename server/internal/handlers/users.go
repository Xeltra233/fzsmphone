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
		SELECT id, discord_id, username, display_name, COALESCE(avatar_url, ''), role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query users")
		return
	}
	defer rows.Close()

	type userResp struct {
		ID          int64     `json:"id"`
		DiscordID   string    `json:"discord_id"`
		Username    string    `json:"username"`
		DisplayName string    `json:"display_name"`
		AvatarURL   string    `json:"avatar_url"`
		Role        string    `json:"role"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	var users []userResp
	for rows.Next() {
		var u userResp
		if err := rows.Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName,
			&u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
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
		ID          int64     `json:"id"`
		DiscordID   string    `json:"discord_id"`
		Username    string    `json:"username"`
		DisplayName string    `json:"display_name"`
		AvatarURL   string    `json:"avatar_url"`
		Role        string    `json:"role"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, discord_id, username, display_name, COALESCE(avatar_url, ''), role, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName,
		&u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		mw.Error(w, http.StatusNotFound, "user not found")
		return
	}

	mw.JSON(w, http.StatusOK, u)
}

// PATCH /api/users/{id}/role
func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Only admins can update roles - verify caller is admin
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

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Role != "admin" && body.Role != "user" && body.Role != "moderator" {
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
