package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type CreditsHandler struct {
	DB *database.DB
}

type SigninResponse struct {
	CreditsEarned int `json:"credits_earned"`
	StreakBonus   int `json:"streak_bonus"`
	TotalCredits  int `json:"total_credits"`
	Streak        int `json:"streak"`
}

func (h *CreditsHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var lastSignin *time.Time
	var streak int
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT last_signin_at, signin_streak FROM users WHERE id = $1
	`, userID).Scan(&lastSignin, &streak)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to check signin")
		return
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	if lastSignin != nil {
		lastDay := lastSignin.Truncate(24 * time.Hour)
		if lastDay.Equal(today) {
			mw.Error(w, http.StatusBadRequest, "今日已签到")
			return
		}
		if lastDay.Add(24*time.Hour).Equal(today) || lastDay.Add(48*time.Hour).Equal(today) {
			streak++
		} else {
			streak = 1
		}
	} else {
		streak = 1
	}

	var dailyCredits, streakBonus int
	h.DB.Pool.QueryRow(r.Context(), `
		SELECT COALESCE(value, '10')::int, COALESCE((SELECT value FROM app_settings WHERE key = 'signin_streak_bonus'), '5')::int
		FROM app_settings WHERE key = 'signin_daily_credits'
	`).Scan(&dailyCredits, &streakBonus)

	creditsEarned := dailyCredits
	if streak > 1 {
		creditsEarned += streakBonus * (streak - 1)
	}

	_, err = h.DB.Pool.Exec(r.Context(), `
		UPDATE users SET 
			credits = credits + $1,
			last_signin_at = NOW(),
			signin_streak = $2
		WHERE id = $3
	`, creditsEarned, streak, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "签到失败")
		return
	}

	h.DB.Pool.Exec(r.Context(), `
		INSERT INTO signin_records (user_id, credits_earned, streak_bonus)
		VALUES ($1, $2, $3)
	`, userID, dailyCredits, creditsEarned-dailyCredits)

	var totalCredits int
	h.DB.Pool.QueryRow(r.Context(), `SELECT credits FROM users WHERE id = $1`, userID).Scan(&totalCredits)

	mw.JSON(w, http.StatusOK, SigninResponse{
		CreditsEarned: creditsEarned,
		StreakBonus:   creditsEarned - dailyCredits,
		TotalCredits:  totalCredits,
		Streak:        streak,
	})
}

func (h *CreditsHandler) GetCredits(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var credits, totalTokens, streak int
	var lastSignin *time.Time
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT credits, total_tokens, signin_streak, last_signin_at FROM users WHERE id = $1
	`, userID).Scan(&credits, &totalTokens, &streak, &lastSignin)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get credits")
		return
	}

	canSignIn := true
	if lastSignin != nil {
		today := time.Now().Truncate(24 * time.Hour)
		lastDay := lastSignin.Truncate(24 * time.Hour)
		if lastDay.Equal(today) {
			canSignIn = false
		}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"credits":      credits,
		"total_tokens": totalTokens,
		"streak":       streak,
		"can_signin":   canSignIn,
	})
}

func (h *CreditsHandler) GetInviteCode(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var inviteCode string
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT COALESCE(invite_code, '') FROM users WHERE id = $1
	`, userID).Scan(&inviteCode)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get invite code")
		return
	}

	if inviteCode == "" {
		inviteCode = generateInviteCode()
		h.DB.Pool.Exec(r.Context(), `UPDATE users SET invite_code = $1 WHERE id = $2`, inviteCode, userID)
	}

	var inviteCount, rewardsClaimed int
	h.DB.Pool.QueryRow(r.Context(), `
		SELECT COUNT(*), COALESCE(SUM(reward_credits), 0) FROM invite_rewards WHERE inviter_id = $1
	`, userID).Scan(&inviteCount, &rewardsClaimed)

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"invite_code":     inviteCode,
		"invite_count":    inviteCount,
		"rewards_claimed": rewardsClaimed,
	})
}

func (h *CreditsHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var credits, totalTokens, streak, inviteRewards int
	var lastSignin *time.Time
	var inviteCode, modelDisplayNames string

	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT credits, total_tokens, signin_streak, last_signin_at, COALESCE(invite_code, ''), COALESCE(model_display_names::text, '{}'), COALESCE(invite_rewards_claimed, 0)
		FROM users WHERE id = $1
	`, userID).Scan(&credits, &totalTokens, &streak, &lastSignin, &inviteCode, &modelDisplayNames, &inviteRewards)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get settings")
		return
	}

	var defaultCredits, signinCredits, streakBonus, inviteCredits int
	var signinEnabled, inviteEnabled bool
	h.DB.Pool.QueryRow(r.Context(), `
		SELECT 
			COALESCE((SELECT value::int FROM app_settings WHERE key = 'default_credits'), 1000),
			COALESCE((SELECT value::int FROM app_settings WHERE key = 'signin_daily_credits'), 10),
			COALESCE((SELECT value::int FROM app_settings WHERE key = 'signin_streak_bonus'), 5),
			COALESCE((SELECT value::int FROM app_settings WHERE key = 'invite_reward_credits'), 100),
			COALESCE((SELECT value::bool FROM app_settings WHERE key = 'signin_enabled'), true),
			COALESCE((SELECT value::bool FROM app_settings WHERE key = 'invite_enabled'), true)
	`).Scan(&defaultCredits, &signinCredits, &streakBonus, &inviteCredits, &signinEnabled, &inviteEnabled)

	canSignIn := true
	if lastSignin != nil {
		today := time.Now().Truncate(24 * time.Hour)
		lastDay := lastSignin.Truncate(24 * time.Hour)
		if lastDay.Equal(today) {
			canSignIn = false
		}
	}

	var displayNames map[string]string
	json.Unmarshal([]byte(modelDisplayNames), &displayNames)

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"credits":         credits,
		"total_tokens":    totalTokens,
		"streak":          streak,
		"can_signin":      canSignIn,
		"invite_code":     inviteCode,
		"invite_count":    0,
		"rewards_claimed": inviteRewards,
		"settings": map[string]interface{}{
			"default_credits":       defaultCredits,
			"signin_daily_credits":  signinCredits,
			"signin_streak_bonus":   streakBonus,
			"invite_reward_credits": inviteCredits,
			"signin_enabled":        signinEnabled,
			"invite_enabled":        inviteEnabled,
		},
		"model_display_names": displayNames,
	})
}

func (h *CreditsHandler) UpdateModelDisplayName(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Model       string `json:"model"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	var currentNames map[string]string
	h.DB.Pool.QueryRow(r.Context(), `
		SELECT COALESCE(model_display_names::text, '{}') FROM users WHERE id = $1
	`, userID).Scan(&currentNames)

	if currentNames == nil {
		currentNames = make(map[string]string)
	}

	currentNames[body.Model] = body.DisplayName

	namesJSON, _ := json.Marshal(currentNames)
	h.DB.Pool.Exec(r.Context(), `
		UPDATE users SET model_display_names = $1 WHERE id = $2
	`, string(namesJSON), userID)

	mw.JSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *CreditsHandler) ConsumeCredits(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Tokens int    `json:"tokens"`
		Model  string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	// 1000 tokens = 1 credit
	cost := (body.Tokens + 999) / 1000
	if cost < 1 {
		cost = 1
	}

	var credits int
	err := h.DB.Pool.QueryRow(r.Context(), `
		UPDATE users SET 
			credits = credits - $1,
			total_tokens = total_tokens + $2
		WHERE id = $3
		AND credits >= $1
		RETURNING credits
	`, cost, body.Tokens, userID).Scan(&credits)

	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusForbidden, "额度不足")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "消费失败")
		return
	}

	h.DB.Pool.Exec(r.Context(), `
		INSERT INTO usage_records (user_id, tokens_used, model, credits_cost)
		VALUES ($1, $2, $3, $4)
	`, userID, body.Tokens, body.Model, cost)

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"cost":    cost,
		"credits": credits,
		"tokens":  body.Tokens,
	})
}

func (h *CreditsHandler) UseInviteCode(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	// Check if user already used an invite code
	var invitedBy *int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT invited_by FROM users WHERE id = $1
	`, userID).Scan(&invitedBy)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to check")
		return
	}
	if invitedBy != nil {
		mw.Error(w, http.StatusBadRequest, "已使用过邀请码")
		return
	}

	var body struct {
		InviteCode string `json:"invite_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Find inviter
	var inviterID int64
	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id FROM users WHERE invite_code = $1 AND id != $2
	`, body.InviteCode, userID).Scan(&inviterID)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusBadRequest, "邀请码无效")
		return
	}

	// Get invite reward amount
	var rewardCredits int
	h.DB.Pool.QueryRow(r.Context(), `
		SELECT COALESCE(value::int, 100) FROM app_settings WHERE key = 'invite_reward_credits'
	`).Scan(&rewardCredits)

	// Update invitee
	_, err = h.DB.Pool.Exec(r.Context(), `
		UPDATE users SET invited_by = $1, credits = credits + $2 WHERE id = $3
	`, inviterID, rewardCredits, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "绑定邀请码失败")
		return
	}

	// Give reward to inviter
	_, err = h.DB.Pool.Exec(r.Context(), `
		INSERT INTO invite_rewards (inviter_id, invited_id, reward_credits)
		VALUES ($1, $2, $3)
	`, inviterID, userID, rewardCredits)
	if err == nil {
		h.DB.Pool.Exec(r.Context(), `
			UPDATE users SET credits = credits + $1, invite_rewards_claimed = invite_rewards_claimed + $1 WHERE id = $2
		`, rewardCredits, inviterID)
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"message":        "绑定成功",
		"credits_earned": rewardCredits,
	})
}

func generateInviteCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}

func (h *CreditsHandler) RegisterRoutes(r chi.Router) {
	r.Post("/signin", h.SignIn)
	r.Get("/credits", h.GetCredits)
	r.Get("/invite-code", h.GetInviteCode)
	r.Post("/invite-code", h.UseInviteCode)
	r.Post("/consume", h.ConsumeCredits)
	r.Get("/settings", h.GetUserSettings)
	r.Put("/model-display-name", h.UpdateModelDisplayName)
}
