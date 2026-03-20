package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type CreditsHandler struct {
	DB     *database.DB
	Logger *LoggerHandler
}

type SigninResponse struct {
	CreditsEarned int `json:"credits_earned"`
	StreakBonus   int `json:"streak_bonus"`
	TotalCredits  int `json:"total_credits"`
	Streak        int `json:"streak"`
}

func (h *CreditsHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var signinEnabled bool
	if err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT COALESCE((SELECT value::bool FROM app_settings WHERE key = 'signin_enabled'), true)
	`).Scan(&signinEnabled); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to load signin settings")
		return
	}
	if !signinEnabled {
		mw.Error(w, http.StatusBadRequest, "签到功能未开启")
		return
	}

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

	// Give 15% sign-in rebate to inviter if exists
	var inviterID int64
	err = h.DB.Pool.QueryRow(r.Context(), `SELECT COALESCE(invited_by, 0) FROM users WHERE id = $1`, userID).Scan(&inviterID)
	if err == nil && inviterID > 0 {
		rebate := float64(creditsEarned) * 0.15
		if rebate >= 1 {
			rebateCredits := int(rebate)
			h.DB.Pool.Exec(r.Context(), `UPDATE users SET credits = credits + $1 WHERE id = $2`, rebateCredits, inviterID)
			h.DB.Pool.Exec(r.Context(), `INSERT INTO invite_rewards (inviter_id, invited_id, reward_credits, reward_type) VALUES ($1, $2, $3, 'signin_rebate')`, inviterID, userID, rebateCredits)
		}
	}

	h.DB.Pool.Exec(r.Context(), `
	INSERT INTO signin_records (user_id, credits_earned, streak_bonus)
	VALUES ($1, $2, $3)
	`, userID, dailyCredits, creditsEarned-dailyCredits)

	var totalCredits int
	h.DB.Pool.QueryRow(r.Context(), `SELECT credits FROM users WHERE id = $1`, userID).Scan(&totalCredits)

	if h.Logger != nil {
		var username string
		h.DB.Pool.QueryRow(r.Context(), `SELECT username FROM users WHERE id = $1`, userID).Scan(&username)
		h.Logger.Log("INFO", "user_signin", &userID, &username, map[string]interface{}{
			"credits_earned": creditsEarned,
			"streak":         streak,
		}, r.RemoteAddr)
	}

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
	var signinEnabled bool
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT credits, total_tokens, signin_streak, last_signin_at FROM users WHERE id = $1
	`, userID).Scan(&credits, &totalTokens, &streak, &lastSignin)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get credits")
		return
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT COALESCE((SELECT value::bool FROM app_settings WHERE key = 'signin_enabled'), true)
	`).Scan(&signinEnabled)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get signin settings")
		return
	}

	canSignIn := signinEnabled
	if signinEnabled && lastSignin != nil {
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
	var announcement, tips string
	h.DB.Pool.QueryRow(r.Context(), `
	SELECT
	COALESCE((SELECT value::int FROM app_settings WHERE key = 'default_credits'), 1000),
	COALESCE((SELECT value::int FROM app_settings WHERE key = 'signin_daily_credits'), 10),
	COALESCE((SELECT value::int FROM app_settings WHERE key = 'signin_streak_bonus'), 5),
	COALESCE((SELECT value::int FROM app_settings WHERE key = 'invite_reward_credits'), 100),
	COALESCE((SELECT value::bool FROM app_settings WHERE key = 'signin_enabled'), true),
	COALESCE((SELECT value::bool FROM app_settings WHERE key = 'invite_enabled'), true),
	COALESCE((SELECT value FROM app_settings WHERE key = 'announcement'), ''),
	COALESCE((SELECT value FROM app_settings WHERE key = 'tips'), '')
	`).Scan(&defaultCredits, &signinCredits, &streakBonus, &inviteCredits, &signinEnabled, &inviteEnabled, &announcement, &tips)

	canSignIn := signinEnabled
	if signinEnabled && lastSignin != nil {
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
		"announcement":    announcement,
		"tips":            tips,
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

func (h *CreditsHandler) GetInviteInfo(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	type invitee struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Reward   int    `json:"reward"`
	}

	var code string
	var totalRewards int
	var invitees []invitee

	if err := h.DB.Pool.QueryRow(r.Context(), `SELECT COALESCE(invite_code, '') FROM users WHERE id = $1`, userID).Scan(&code); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get invite info")
		return
	}

	if code == "" {
		code = generateInviteCode()
		if _, err := h.DB.Pool.Exec(r.Context(), `UPDATE users SET invite_code = $1 WHERE id = $2`, code, userID); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to create invite code")
			return
		}
	}

	if err := h.DB.Pool.QueryRow(r.Context(), `SELECT COALESCE(SUM(reward_credits), 0) FROM invite_rewards WHERE inviter_id = $1`, userID).Scan(&totalRewards); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to load invite rewards")
		return
	}

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT u.id, u.username, COALESCE(SUM(r.reward_credits), 0) as total_reward
		FROM users u
		LEFT JOIN invite_rewards r ON r.invited_id = u.id
		WHERE u.invited_by = $1
		GROUP BY u.id, u.username
		ORDER BY u.created_at DESC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to load invite list")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var inv invitee
		if err := rows.Scan(&inv.ID, &inv.Username, &inv.Reward); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan invite list")
			return
		}
		invitees = append(invitees, inv)
	}
	if err := rows.Err(); err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to read invite list")
		return
	}
	if invitees == nil {
		invitees = []invitee{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"code":         code,
		"invitees":     invitees,
		"totalRewards": totalRewards,
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

func (h *CreditsHandler) RedeemCoupon(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	if body.Code == "" {
		mw.Error(w, http.StatusBadRequest, "请输入兑换码")
		return
	}

	// Find coupon
	var coupon struct {
		ID          int64      `json:"id"`
		Credits     int        `json:"credits"`
		MaxUses     int        `json:"max_uses"`
		CurrentUses int        `json:"current_uses"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}
	err := h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, credits, max_uses, current_uses, expires_at 
		FROM coupon_codes 
		WHERE code = $1
	`, body.Code).Scan(&coupon.ID, &coupon.Credits, &coupon.MaxUses, &coupon.CurrentUses, &coupon.ExpiresAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusBadRequest, "兑换码无效")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询失败")
		return
	}

	// Check if expired
	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
		mw.Error(w, http.StatusBadRequest, "兑换码已过期")
		return
	}

	// Check usage limit
	if coupon.CurrentUses >= coupon.MaxUses {
		mw.Error(w, http.StatusBadRequest, "兑换码已使用完毕")
		return
	}

	// Check if user already used this coupon
	var alreadyUsed bool
	h.DB.Pool.QueryRow(r.Context(), `
		SELECT EXISTS(SELECT 1 FROM coupon_redeem_records WHERE user_id = $1 AND coupon_id = $2)
	`, userID, coupon.ID).Scan(&alreadyUsed)
	if alreadyUsed {
		mw.Error(w, http.StatusBadRequest, "您已使用过此兑换码")
		return
	}

	// Redeem coupon
	tx, err := h.DB.Pool.Begin(r.Context())
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "兑换失败")
		return
	}
	defer tx.Rollback(r.Context())

	// Update coupon usage
	_, err = tx.Exec(r.Context(), `
		UPDATE coupon_codes SET current_uses = current_uses + 1 WHERE id = $1
	`, coupon.ID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "兑换失败")
		return
	}

	// Add credits to user
	_, err = tx.Exec(r.Context(), `
		UPDATE users SET credits = credits + $1 WHERE id = $2
	`, coupon.Credits, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "添加额度失败")
		return
	}

	// Record redemption
	_, err = tx.Exec(r.Context(), `
		INSERT INTO coupon_redeem_records (user_id, coupon_id, credits_earned)
		VALUES ($1, $2, $3)
	`, userID, coupon.ID, coupon.Credits)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "记录失败")
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "兑换失败")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"message":        "兑换成功",
		"credits_earned": coupon.Credits,
	})
}

// Admin: Create coupon code
func (h *CreditsHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "需要超级管理员权限")
		return
	}

	var body struct {
		Code      string `json:"code"`
		Credits   int    `json:"credits"`
		MaxUses   int    `json:"max_uses"`
		ExpiresIn int    `json:"expires_in"` // days, 0 = never
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	if body.Code == "" {
		body.Code = generateInviteCode()
	}
	if body.Credits <= 0 {
		mw.Error(w, http.StatusBadRequest, "额度必须大于0")
		return
	}
	if body.MaxUses <= 0 {
		body.MaxUses = 1
	}

	var expiresAt *time.Time
	if body.ExpiresIn > 0 {
		t := time.Now().AddDate(0, 0, body.ExpiresIn)
		expiresAt = &t
	}

	_, err := h.DB.Pool.Exec(r.Context(), `
		INSERT INTO coupon_codes (code, credits, max_uses, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (code) DO UPDATE SET credits = $2, max_uses = $3, expires_at = $4
	`, body.Code, body.Credits, body.MaxUses, expiresAt)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "创建失败")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"code":     body.Code,
		"credits":  body.Credits,
		"max_uses": body.MaxUses,
	})
}

// Admin: List coupon codes
func (h *CreditsHandler) ListCoupons(w http.ResponseWriter, r *http.Request) {
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "需要超级管理员权限")
		return
	}

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT code, credits, max_uses, current_uses, expires_at, created_at
		FROM coupon_codes
		ORDER BY created_at DESC
		LIMIT 50
	`)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询失败")
		return
	}
	defer rows.Close()

	type couponRow struct {
		Code        string     `json:"code"`
		Credits     int        `json:"credits"`
		MaxUses     int        `json:"max_uses"`
		CurrentUses int        `json:"current_uses"`
		ExpiresAt   *time.Time `json:"expires_at"`
		CreatedAt   time.Time  `json:"created_at"`
	}

	var coupons []couponRow
	for rows.Next() {
		var c couponRow
		rows.Scan(&c.Code, &c.Credits, &c.MaxUses, &c.CurrentUses, &c.ExpiresAt, &c.CreatedAt)
		coupons = append(coupons, c)
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"coupons": coupons,
	})
}

func (h *CreditsHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT username, COALESCE(total_tokens, 0) as tokens FROM users 
		WHERE discord_id IS NOT NULL OR email IS NOT NULL
		ORDER BY tokens DESC LIMIT 10
	`)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get leaderboard")
		return
	}
	defer rows.Close()

	var leaderboard []map[string]interface{}
	rank := 1
	for rows.Next() {
		var username string
		var tokens int
		if err := rows.Scan(&username, &tokens); err != nil {
			continue
		}
		level := calculateLevel(tokens)
		leaderboard = append(leaderboard, map[string]interface{}{
			"name":  username,
			"score": tokens,
			"level": level,
			"rank":  rank,
		})
		rank++
	}

	fakeUsers := []struct {
		name   string
		tokens int
	}{
		{"小明", 152000},
		{"小红", 128000},
		{"小华", 96000},
		{"小李", 74000},
		{"小张", 51000},
		{"小刚", 35000},
		{"小丽", 28000},
		{"小芳", 15000},
		{"小军", 9000},
		{"小美", 5000},
	}

	for _, fake := range fakeUsers {
		exists := false
		for _, l := range leaderboard {
			if l["name"] == fake.name {
				exists = true
				break
			}
		}
		if !exists {
			level := calculateLevel(fake.tokens)
			leaderboard = append(leaderboard, map[string]interface{}{
				"name":  fake.name,
				"score": fake.tokens,
				"level": level,
			})
		}
	}

	if len(leaderboard) < 5 {
		needed := 5 - len(leaderboard)
		for i := 0; i < needed && i < len(fakeUsers); i++ {
			level := calculateLevel(fakeUsers[i].tokens)
			leaderboard = append(leaderboard, map[string]interface{}{
				"name":  fakeUsers[i].name,
				"score": fakeUsers[i].tokens,
				"level": level,
			})
		}
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i]["score"].(int) > leaderboard[j]["score"].(int)
	})

	for i := range leaderboard {
		leaderboard[i]["rank"] = i + 1
	}

	mw.JSON(w, http.StatusOK, leaderboard)
}

func calculateLevel(tokens int) string {
	switch {
	case tokens >= 1000000:
		return "传奇"
	case tokens >= 500000:
		return "王者"
	case tokens >= 200000:
		return "钻石"
	case tokens >= 100000:
		return "铂金"
	case tokens >= 50000:
		return "黄金"
	case tokens >= 20000:
		return "白银"
	case tokens >= 5000:
		return "青铜"
	default:
		return "入门"
	}
}

func (h *CreditsHandler) RegisterRoutes(r chi.Router) {
	r.Post("/signin", h.SignIn)
	r.Get("/credits", h.GetCredits)
	r.Get("/leaderboard", h.GetLeaderboard)
	r.Get("/invite-code", h.GetInviteCode)
	r.Post("/invite-code", h.UseInviteCode)
	r.Get("/invite-info", h.GetInviteInfo)
	r.Post("/consume", h.ConsumeCredits)
	r.Get("/settings", h.GetUserSettings)
	r.Put("/model-display-name", h.UpdateModelDisplayName)
	r.Post("/redeem", h.RedeemCoupon)
	r.Post("/coupon", h.CreateCoupon)
	r.Get("/coupons", h.ListCoupons)
}
