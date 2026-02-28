package handlers

import (
	"encoding/json"
	"net/http"

	"fzsmphone/internal/database"
	"fzsmphone/internal/middleware"
)

type GameHandler struct {
	DB *database.DB
}

// ListRecords returns the last 50 game records.
func (h *GameHandler) ListRecords(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(),
		`SELECT id, game, detail, amount, win, created_at
		 FROM game_records
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT 50`, userID)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	type Record struct {
		ID        int64   `json:"id"`
		Game      string  `json:"game"`
		Detail    string  `json:"detail"`
		Amount    float64 `json:"amount"`
		Win       bool    `json:"win"`
		CreatedAt string  `json:"created_at"`
	}

	var records []Record
	for rows.Next() {
		var rec Record
		var t interface{}
		if err := rows.Scan(&rec.ID, &rec.Game, &rec.Detail, &rec.Amount, &rec.Win, &t); err != nil {
			continue
		}
		if ct, ok := t.(interface{ String() string }); ok {
			rec.CreatedAt = ct.String()
		}
		records = append(records, rec)
	}
	if records == nil {
		records = []Record{}
	}
	middleware.JSON(w, http.StatusOK, map[string]interface{}{"data": records})
}

// CreateRecord adds a new game record.
func (h *GameHandler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	var input struct {
		Game   string  `json:"game"`
		Detail string  `json:"detail"`
		Amount float64 `json:"amount"`
		Win    bool    `json:"win"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if input.Game == "" {
		middleware.Error(w, http.StatusBadRequest, "game is required")
		return
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(),
		`INSERT INTO game_records (user_id, game, detail, amount, win)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		userID, input.Game, input.Detail, input.Amount, input.Win,
	).Scan(&id)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, "insert failed")
		return
	}

	middleware.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "game record created",
	})
}

// ClearRecords deletes all game records for the user.
func (h *GameHandler) ClearRecords(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	_, err := h.DB.Pool.Exec(r.Context(),
		`DELETE FROM game_records WHERE user_id = $1`, userID)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, "delete failed")
		return
	}

	middleware.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "game records cleared",
	})
}
