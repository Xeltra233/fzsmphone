package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/jackc/pgx/v5"
)

type WalletHandler struct {
	DB *database.DB
}

// GET /api/wallet
func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	// Ensure wallet exists (auto-create with default balance)
	_, err := h.DB.Pool.Exec(r.Context(), `
		INSERT INTO wallets (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to ensure wallet")
		return
	}

	var balance float64
	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT balance FROM wallets WHERE user_id = $1
	`, userID).Scan(&balance)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get wallet")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"balance": balance,
	})
}

// GET /api/wallet/transactions
func (h *WalletHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT id, user_id, type, amount, description, target, created_at
		FROM wallet_transactions WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 200
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query transactions")
		return
	}
	defer rows.Close()

	type txResp struct {
		ID          int64   `json:"id"`
		UserID      int64   `json:"user_id"`
		Type        string  `json:"type"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Target      string  `json:"target"`
		CreatedAt   string  `json:"created_at"`
	}

	var txs []txResp
	for rows.Next() {
		var t txResp
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &t.Amount,
			&t.Description, &t.Target, &t.CreatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan transaction")
			return
		}
		txs = append(txs, t)
	}
	if txs == nil {
		txs = []txResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": txs})
}

// POST /api/wallet/transaction
func (h *WalletHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Type        string  `json:"type"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Target      string  `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Ensure wallet exists
	_, err := h.DB.Pool.Exec(r.Context(), `
		INSERT INTO wallets (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to ensure wallet")
		return
	}

	// Check balance for debit transactions
	if body.Amount < 0 {
		var balance float64
		err = h.DB.Pool.QueryRow(r.Context(), `SELECT balance FROM wallets WHERE user_id = $1`, userID).Scan(&balance)
		if err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to check balance")
			return
		}
		if balance+body.Amount < 0 {
			mw.Error(w, http.StatusBadRequest, "insufficient balance")
			return
		}
	}

	// Update balance
	var newBalance float64
	err = h.DB.Pool.QueryRow(r.Context(), `
		UPDATE wallets SET balance = balance + $1 WHERE user_id = $2 RETURNING balance
	`, body.Amount, userID).Scan(&newBalance)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "wallet not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update balance")
		return
	}

	// Record transaction
	var txID int64
	err = h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO wallet_transactions (user_id, type, amount, description, target)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, body.Type, math.Abs(body.Amount), body.Description, body.Target).Scan(&txID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to record transaction")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":      txID,
		"balance": newBalance,
		"message": "transaction recorded",
	})
}

// PUT /api/wallet/balance
func (h *WalletHandler) SetBalance(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Balance float64 `json:"balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, err := h.DB.Pool.Exec(r.Context(), `
		INSERT INTO wallets (user_id, balance) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET balance = $2
	`, userID, body.Balance)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to set balance")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"balance": body.Balance,
		"message": "balance updated",
	})
}
