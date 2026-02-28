package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type CallHandler struct {
	DB *database.DB
}

// GET /api/calls
func (h *CallHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, character_id, name, number, avatar, type, call_type, duration, created_at
		FROM call_records WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query calls")
		return
	}
	defer rows.Close()

	type callResp struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		CharacterID string `json:"character_id"`
		Name        string `json:"name"`
		Number      string `json:"number"`
		Avatar      string `json:"avatar"`
		Type        string `json:"type"`
		CallType    string `json:"call_type"`
		Duration    string `json:"duration"`
		CreatedAt   string `json:"created_at"`
	}

	var calls []callResp
	for rows.Next() {
		var c callResp
		if err := rows.Scan(&c.ID, &c.UserID, &c.CharacterID, &c.Name, &c.Number,
			&c.Avatar, &c.Type, &c.CallType, &c.Duration, &c.CreatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan call")
			return
		}
		calls = append(calls, c)
	}
	if calls == nil {
		calls = []callResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": calls})
}

// POST /api/calls
func (h *CallHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		CharacterID string `json:"character_id"`
		Name        string `json:"name"`
		Number      string `json:"number"`
		Avatar      string `json:"avatar"`
		Type        string `json:"type"`
		CallType    string `json:"call_type"`
		Duration    string `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO call_records (user_id, character_id, name, number, avatar, type, call_type, duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, userID, body.CharacterID, body.Name, body.Number, body.Avatar,
		body.Type, body.CallType, body.Duration).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create call record")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "call record created"})
}

// DELETE /api/calls/{id}
func (h *CallHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM call_records WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete call record")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "call record not found")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "call record deleted"})
}

// DELETE /api/calls
func (h *CallHandler) Clear(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	_, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM call_records WHERE user_id = $1
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to clear call records")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "call records cleared"})
}
