package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type LoggerHandler struct {
	DB *database.DB
}

type LogEntry struct {
	ID        int64                  `json:"id"`
	Level     string                 `json:"level"`
	Action    string                 `json:"action"`
	UserID    *int64                 `json:"user_id,omitempty"`
	Username  *string                `json:"username,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	IPAddress *string                `json:"ip_address,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

func (h *LoggerHandler) Log(level, action string, userID *int64, username *string, details map[string]interface{}, ip string) {
	if details == nil {
		details = make(map[string]interface{})
	}

	h.DB.Pool.Exec(context.Background(), `
		INSERT INTO system_logs (level, action, user_id, username, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, level, action, userID, username, details, ip)
}

func (h *LoggerHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "需要超级管理员权限")
		return
	}

	action := r.URL.Query().Get("action")
	level := r.URL.Query().Get("level")
	userIDStr := r.URL.Query().Get("user_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o > 0 {
		offset = o
	}

	query := `
		SELECT id, level, action, user_id, username, details, ip_address, created_at
		FROM system_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if action != "" {
		query += ` AND action = $` + strconv.Itoa(argNum)
		args = append(args, action)
		argNum++
	}
	if level != "" {
		query += ` AND level = $` + strconv.Itoa(argNum)
		args = append(args, level)
		argNum++
	}
	if userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			query += ` AND user_id = $` + strconv.Itoa(argNum)
			args = append(args, userID)
			argNum++
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argNum) + ` OFFSET $` + strconv.Itoa(argNum+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Pool.Query(r.Context(), query, args...)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询失败")
		return
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		var details []byte
		var username *string
		var userID *int64
		var ip *string

		err := rows.Scan(&log.ID, &log.Level, &log.Action, &userID, &username, &details, &ip, &log.CreatedAt)
		if err != nil {
			continue
		}

		log.UserID = userID
		log.Username = username
		log.IPAddress = ip
		if details != nil {
			json.Unmarshal(details, &log.Details)
		}

		logs = append(logs, log)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM system_logs WHERE 1=1`
	countArgs := []interface{}{}
	countNum := 1

	if action != "" {
		countQuery += ` AND action = $` + strconv.Itoa(countNum)
		countArgs = append(countArgs, action)
		countNum++
	}
	if level != "" {
		countQuery += ` AND level = $` + strconv.Itoa(countNum)
		countArgs = append(countArgs, level)
		countNum++
	}
	if userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			countQuery += ` AND user_id = $` + strconv.Itoa(countNum)
			countArgs = append(countArgs, userID)
		}
	}

	h.DB.Pool.QueryRow(r.Context(), countQuery, countArgs...).Scan(&total)

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"logs":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *LoggerHandler) GetLogStats(w http.ResponseWriter, r *http.Request) {
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "需要超级管理员权限")
		return
	}

	var stats struct {
		TodayCount     int `json:"today_count"`
		YesterdayCount int `json:"yesterday_count"`
		WeekCount      int `json:"week_count"`
		ErrorCount     int `json:"error_count"`
	}

	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	weekAgo := today.AddDate(0, 0, -7)

	h.DB.Pool.QueryRow(r.Context(), `
		SELECT 
			(SELECT COUNT(*) FROM system_logs WHERE created_at >= $1) as today,
			(SELECT COUNT(*) FROM system_logs WHERE created_at >= $2 AND created_at < $1) as yesterday,
			(SELECT COUNT(*) FROM system_logs WHERE created_at >= $3) as week,
			(SELECT COUNT(*) FROM system_logs WHERE level = 'ERROR' AND created_at >= $1) as errors
	`, today, yesterday, weekAgo).Scan(&stats.TodayCount, &stats.YesterdayCount, &stats.WeekCount, &stats.ErrorCount)

	// Get action counts
	rows, _ := h.DB.Pool.Query(r.Context(), `
		SELECT action, COUNT(*) as cnt FROM system_logs 
		WHERE created_at >= $1 
		GROUP BY action 
		ORDER BY cnt DESC 
		LIMIT 10
	`, weekAgo)

	var actionCounts []map[string]interface{}
	for rows.Next() {
		var action string
		var count int
		rows.Scan(&action, &count)
		actionCounts = append(actionCounts, map[string]interface{}{
			"action": action,
			"count":  count,
		})
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"stats":         stats,
		"action_counts": actionCounts,
	})
}

func (h *LoggerHandler) ClearLogs(w http.ResponseWriter, r *http.Request) {
	isSuperAdmin, _ := mw.GetIsSuperAdmin(r.Context())
	if !isSuperAdmin {
		mw.Error(w, http.StatusForbidden, "需要超级管理员权限")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	result, err := h.DB.Pool.Exec(r.Context(), `
	DELETE FROM system_logs WHERE created_at < $1
	`, cutoff)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "清除失败")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{
		"message":      "清除成功",
		"deleted_rows": result.RowsAffected(),
	})
}

func (h *LoggerHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.ListLogs)
	r.Get("/stats", h.GetLogStats)
	r.Delete("/", h.ClearLogs)
}

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxy)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}
	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Middleware to log HTTP requests
func (h *LoggerHandler) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get user info if authenticated
		var userID *int64
		var username *string
		if uid, ok := mw.GetUserID(r.Context()); ok {
			u := uid
			userID = &u
			if un, ok := mw.GetUserRole(r.Context()); ok {
				username = &un
			}
		}

		// Call next handler
		next.ServeHTTP(w, r)

		// Log the request
		duration := time.Since(start)
		details := map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.Milliseconds(),
			"status":   200, // Simplified
		}

		h.Log("INFO", "http_request", userID, username, details, getClientIP(r))
	})
}
