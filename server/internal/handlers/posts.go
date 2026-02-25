package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fzsmphone/internal/database"
	mw "fzsmphone/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type PostHandler struct {
	DB *database.DB
}

// GET /api/posts
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	board := r.URL.Query().Get("board")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	query := `
		SELECT p.id, p.user_id, p.board, p.title, p.content, p.images, p.likes,
		       p.comments, p.pinned, p.created_at, p.updated_at,
		       u.username, u.avatar_url
		FROM posts p
		LEFT JOIN users u ON u.id = p.user_id
	`
	args := []interface{}{}
	argIdx := 1

	if board != "" {
		query += ` WHERE p.board = $` + strconv.Itoa(argIdx)
		args = append(args, board)
		argIdx++
	}

	query += ` ORDER BY p.pinned DESC, p.created_at DESC`
	query += ` LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Pool.Query(r.Context(), query, args...)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query posts")
		return
	}
	defer rows.Close()

	type postResp struct {
		ID        int64    `json:"id"`
		UserID    int64    `json:"user_id"`
		Board     string   `json:"board"`
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Images    []string `json:"images"`
		Likes     int      `json:"likes"`
		Comments  int      `json:"comments"`
		Pinned    bool     `json:"pinned"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
		Author    *struct {
			Username  string `json:"username"`
			AvatarURL string `json:"avatar_url"`
		} `json:"author,omitempty"`
	}

	var posts []postResp
	for rows.Next() {
		var p postResp
		var authorName, authorAvatar *string
		if err := rows.Scan(&p.ID, &p.UserID, &p.Board, &p.Title, &p.Content, &p.Images,
			&p.Likes, &p.Comments, &p.Pinned, &p.CreatedAt, &p.UpdatedAt,
			&authorName, &authorAvatar); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan post")
			return
		}
		if p.Images == nil {
			p.Images = []string{}
		}
		if authorName != nil {
			p.Author = &struct {
				Username  string `json:"username"`
				AvatarURL string `json:"avatar_url"`
			}{Username: *authorName}
			if authorAvatar != nil {
				p.Author.AvatarURL = *authorAvatar
			}
		}
		posts = append(posts, p)
	}
	if posts == nil {
		posts = []postResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": posts})
}

// POST /api/posts
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		Board   string   `json:"board"`
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Images  []string `json:"images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Title == "" {
		mw.Error(w, http.StatusBadRequest, "title is required")
		return
	}
	if body.Images == nil {
		body.Images = []string{}
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO posts (user_id, board, title, content, images)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, body.Board, body.Title, body.Content, body.Images).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "post created"})
}

// GET /api/posts/{id}
func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p struct {
		ID        int64    `json:"id"`
		UserID    int64    `json:"user_id"`
		Board     string   `json:"board"`
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Images    []string `json:"images"`
		Likes     int      `json:"likes"`
		Comments  int      `json:"comments"`
		Pinned    bool     `json:"pinned"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
		Username  string   `json:"author_name"`
		AvatarURL string   `json:"author_avatar"`
	}

	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT p.id, p.user_id, p.board, p.title, p.content, p.images, p.likes,
		       p.comments, p.pinned, p.created_at, p.updated_at,
		       COALESCE(u.username, ''), COALESCE(u.avatar_url, '')
		FROM posts p LEFT JOIN users u ON u.id = p.user_id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.UserID, &p.Board, &p.Title, &p.Content, &p.Images,
		&p.Likes, &p.Comments, &p.Pinned, &p.CreatedAt, &p.UpdatedAt,
		&p.Username, &p.AvatarURL)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to get post")
		return
	}
	if p.Images == nil {
		p.Images = []string{}
	}

	mw.JSON(w, http.StatusOK, p)
}

// PUT /api/posts/{id}
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Images  []string `json:"images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Images == nil {
		body.Images = []string{}
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE posts SET title = $1, content = $2, images = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
	`, body.Title, body.Content, body.Images, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to update post")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "post not found or not owned by you")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "post updated"})
}

// DELETE /api/posts/{id}
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	_, _ = h.DB.Pool.Exec(r.Context(), `DELETE FROM comments WHERE post_id = $1`, id)
	result, err := h.DB.Pool.Exec(r.Context(), `
		DELETE FROM posts WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to delete post")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "post not found or not owned by you")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
}

// GET /api/posts/{id}/comments
func (h *PostHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT c.id, c.post_id, c.user_id, c.parent_id, c.content, c.likes, c.created_at,
		       COALESCE(u.username, ''), COALESCE(u.avatar_url, '')
		FROM comments c
		LEFT JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to query comments")
		return
	}
	defer rows.Close()

	type commentResp struct {
		ID        int64  `json:"id"`
		PostID    int64  `json:"post_id"`
		UserID    int64  `json:"user_id"`
		ParentID  *int64 `json:"parent_id,omitempty"`
		Content   string `json:"content"`
		Likes     int    `json:"likes"`
		CreatedAt string `json:"created_at"`
		Username  string `json:"author_name"`
		AvatarURL string `json:"author_avatar"`
	}

	var comments []commentResp
	for rows.Next() {
		var c commentResp
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.ParentID, &c.Content,
			&c.Likes, &c.CreatedAt, &c.Username, &c.AvatarURL); err != nil {
			mw.Error(w, http.StatusInternalServerError, "failed to scan comment")
			return
		}
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []commentResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": comments})
}

// POST /api/posts/{id}/comments
func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Content  string `json:"content"`
		ParentID *int64 `json:"parent_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Content == "" {
		mw.Error(w, http.StatusBadRequest, "content is required")
		return
	}

	var id int64
	err = h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO comments (post_id, user_id, parent_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, postID, userID, body.ParentID, body.Content).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "failed to create comment")
		return
	}

	// Increment comment count
	_, _ = h.DB.Pool.Exec(r.Context(), `
		UPDATE posts SET comments = comments + 1 WHERE id = $1
	`, postID)

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "comment created"})
}
