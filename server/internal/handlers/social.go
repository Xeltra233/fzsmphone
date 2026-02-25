package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// ==================== Weibo ====================

type WeiboHandler struct {
	DB *database.DB
}

// GET /api/weibo
func (h *WeiboHandler) List(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab") // hot, latest, mine
	userID, _ := mw.GetUserID(r.Context())

	query := `
		SELECT w.id, w.user_id, w.content, w.images, w.likes, w.reposts, w.comments,
		       w.is_hot, w.created_at,
		       COALESCE(u.username, ''), COALESCE(u.avatar_url, '')
		FROM weibo_posts w
		LEFT JOIN users u ON u.id = w.user_id
	`
	args := []interface{}{}
	argIdx := 1

	switch tab {
	case "hot":
		query += ` WHERE w.is_hot = true`
	case "mine":
		query += ` WHERE w.user_id = $` + strconv.Itoa(argIdx)
		args = append(args, userID)
		argIdx++
	}

	query += ` ORDER BY w.created_at DESC LIMIT 50`

	rows, err := h.DB.Pool.Query(r.Context(), query, args...)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query weibo")
		return
	}
	defer rows.Close()

	type weiboResp struct {
		ID           int64    `json:"id"`
		UserID       int64    `json:"user_id"`
		Content      string   `json:"content"`
		Images       []string `json:"images"`
		Likes        int      `json:"likes"`
		Reposts      int      `json:"reposts"`
		Comments     int      `json:"comments"`
		IsHot        bool     `json:"is_hot"`
		CreatedAt    string   `json:"created_at"`
		AuthorName   string   `json:"author_name"`
		AuthorAvatar string   `json:"author_avatar"`
	}

	var posts []weiboResp
	for rows.Next() {
		var p weiboResp
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Images, &p.Likes,
			&p.Reposts, &p.Comments, &p.IsHot, &p.CreatedAt,
			&p.AuthorName, &p.AuthorAvatar); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan weibo")
			return
		}
		if p.Images == nil {
			p.Images = []string{}
		}
		posts = append(posts, p)
	}
	if posts == nil {
		posts = []weiboResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": posts})
}

// POST /api/weibo
func (h *WeiboHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Content string   `json:"content"`
		Images  []string `json:"images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Content == "" {
		mw.Error(w, http.StatusBadRequest, "content is required")
		return
	}
	if body.Images == nil {
		body.Images = []string{}
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO weibo_posts (user_id, content, images)
		VALUES ($1, $2, $3)
		RETURNING id
	`, userID, body.Content, body.Images).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create weibo")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "weibo created"})
}

// GET /api/weibo/{id}
func (h *WeiboHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p struct {
		ID           int64    `json:"id"`
		UserID       int64    `json:"user_id"`
		Content      string   `json:"content"`
		Images       []string `json:"images"`
		Likes        int      `json:"likes"`
		Reposts      int      `json:"reposts"`
		Comments     int      `json:"comments"`
		IsHot        bool     `json:"is_hot"`
		CreatedAt    string   `json:"created_at"`
		AuthorName   string   `json:"author_name"`
		AuthorAvatar string   `json:"author_avatar"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT w.id, w.user_id, w.content, w.images, w.likes, w.reposts, w.comments,
		       w.is_hot, w.created_at,
		       COALESCE(u.username, ''), COALESCE(u.avatar_url, '')
		FROM weibo_posts w
		LEFT JOIN users u ON u.id = w.user_id
		WHERE w.id = $1
	`, id).Scan(&p.ID, &p.UserID, &p.Content, &p.Images, &p.Likes,
		&p.Reposts, &p.Comments, &p.IsHot, &p.CreatedAt,
		&p.AuthorName, &p.AuthorAvatar)
	if err != nil {
		mw.Error(w, http.StatusNotFound, "weibo not found")
		return
	}
	if p.Images == nil {
		p.Images = []string{}
	}

	mw.JSON(w, http.StatusOK, p)
}

// DELETE /api/weibo/{id}
func (h *WeiboHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM weibo_posts WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete weibo")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "weibo not found or not owned by you")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "weibo deleted"})
}

// ==================== Moments ====================

type MomentHandler struct {
	DB *database.DB
}

// GET /api/moments
func (h *MomentHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT m.id, m.user_id, m.content, m.images, m.location, m.visibility,
		       m.likes, m.created_at,
		       COALESCE(u.username, ''), COALESCE(u.avatar_url, '')
		FROM moments m
		LEFT JOIN users u ON u.id = m.user_id
		WHERE m.visibility = 'public'
		ORDER BY m.created_at DESC
		LIMIT 50
	`)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query moments")
		return
	}
	defer rows.Close()

	type momentResp struct {
		ID           int64    `json:"id"`
		UserID       int64    `json:"user_id"`
		Content      string   `json:"content"`
		Images       []string `json:"images"`
		Location     string   `json:"location"`
		Visibility   string   `json:"visibility"`
		Likes        int      `json:"likes"`
		CreatedAt    string   `json:"created_at"`
		AuthorName   string   `json:"author_name"`
		AuthorAvatar string   `json:"author_avatar"`
	}

	var moments []momentResp
	for rows.Next() {
		var m momentResp
		if err := rows.Scan(&m.ID, &m.UserID, &m.Content, &m.Images, &m.Location,
			&m.Visibility, &m.Likes, &m.CreatedAt,
			&m.AuthorName, &m.AuthorAvatar); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan moment")
			return
		}
		if m.Images == nil {
			m.Images = []string{}
		}
		moments = append(moments, m)
	}
	if moments == nil {
		moments = []momentResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": moments})
}

// POST /api/moments
func (h *MomentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Content    string   `json:"content"`
		Images     []string `json:"images"`
		Location   string   `json:"location"`
		Visibility string   `json:"visibility"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Images == nil {
		body.Images = []string{}
	}
	if body.Visibility == "" {
		body.Visibility = "public"
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO moments (user_id, content, images, location, visibility)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, body.Content, body.Images, body.Location, body.Visibility).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create moment")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "moment created"})
}

// DELETE /api/moments/{id}
func (h *MomentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM moments WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete moment")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "moment not found or not owned by you")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "moment deleted"})
}
