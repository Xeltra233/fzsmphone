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

type RestaurantHandler struct {
	DB *database.DB
}

type OrderHandler struct {
	DB *database.DB
}

// ==================== Restaurants ====================

// GET /api/restaurants
func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var rows pgx.Rows
	var err error

	if category != "" {
		rows, err = h.DB.Pool.Query(r.Context(), `
			SELECT id, name, category, image_url, rating, delivery_time, min_order, items, created_at
			FROM restaurants WHERE category = $1
			ORDER BY rating DESC
		`, category)
	} else {
		rows, err = h.DB.Pool.Query(r.Context(), `
			SELECT id, name, category, image_url, rating, delivery_time, min_order, items, created_at
			FROM restaurants
			ORDER BY rating DESC
		`)
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询餐厅失败")
		return
	}
	defer rows.Close()

	type restaurantResp struct {
		ID           int64           `json:"id"`
		Name         string          `json:"name"`
		Category     string          `json:"category"`
		ImageURL     string          `json:"image_url"`
		Rating       float64         `json:"rating"`
		DeliveryTime string          `json:"delivery_time"`
		MinOrder     float64         `json:"min_order"`
		Items        json.RawMessage `json:"items"`
		CreatedAt    string          `json:"created_at"`
	}

	var restaurants []restaurantResp
	for rows.Next() {
		var r restaurantResp
		var items []byte
		if err := rows.Scan(&r.ID, &r.Name, &r.Category, &r.ImageURL,
			&r.Rating, &r.DeliveryTime, &r.MinOrder, &items, &r.CreatedAt); err != nil {
			mw.Error(w, http.StatusInternalServerError, "扫描餐厅数据失败")
			return
		}
		if items == nil {
			r.Items = json.RawMessage("[]")
		} else {
			r.Items = json.RawMessage(items)
		}
		restaurants = append(restaurants, r)
	}
	if restaurants == nil {
		restaurants = []restaurantResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": restaurants})
}

// GET /api/restaurants/{id}
func (h *RestaurantHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var rest struct {
		ID           int64           `json:"id"`
		Name         string          `json:"name"`
		Category     string          `json:"category"`
		ImageURL     string          `json:"image_url"`
		Rating       float64         `json:"rating"`
		DeliveryTime string          `json:"delivery_time"`
		MinOrder     float64         `json:"min_order"`
		Items        json.RawMessage `json:"items"`
		CreatedAt    string          `json:"created_at"`
	}

	var items []byte
	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT id, name, category, image_url, rating, delivery_time, min_order, items, created_at
		FROM restaurants WHERE id = $1
	`, id).Scan(&rest.ID, &rest.Name, &rest.Category, &rest.ImageURL,
		&rest.Rating, &rest.DeliveryTime, &rest.MinOrder, &items, &rest.CreatedAt)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "餐厅不存在")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询餐厅失败")
		return
	}
	if items == nil {
		rest.Items = json.RawMessage("[]")
	} else {
		rest.Items = json.RawMessage(items)
	}

	mw.JSON(w, http.StatusOK, rest)
}

// POST /api/restaurants (admin only)
func (h *RestaurantHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	var role string
	_ = h.DB.Pool.QueryRow(r.Context(), `SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	if role != "admin" {
		mw.Error(w, http.StatusForbidden, "需要管理员权限")
		return
	}

	var body struct {
		Name         string          `json:"name"`
		Category     string          `json:"category"`
		ImageURL     string          `json:"image_url"`
		Rating       float64         `json:"rating"`
		DeliveryTime string          `json:"delivery_time"`
		MinOrder     float64         `json:"min_order"`
		Items        json.RawMessage `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "无效的请求体")
		return
	}
	if body.Name == "" {
		mw.Error(w, http.StatusBadRequest, "餐厅名称不能为空")
		return
	}
	if body.Items == nil {
		body.Items = json.RawMessage("[]")
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO restaurants (name, category, image_url, rating, delivery_time, min_order, items)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, body.Name, body.Category, body.ImageURL, body.Rating,
		body.DeliveryTime, body.MinOrder, []byte(body.Items)).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "创建餐厅失败")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "餐厅创建成功"})
}

// ==================== Orders ====================

// GET /api/orders
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	rows, err := h.DB.Pool.Query(r.Context(), `
		SELECT o.id, o.user_id, o.restaurant_id, o.items, o.total, o.status, o.address, o.note,
		       o.created_at, o.updated_at,
		       r.name AS restaurant_name, r.image_url AS restaurant_image
		FROM orders o
		LEFT JOIN restaurants r ON r.id = o.restaurant_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询订单失败")
		return
	}
	defer rows.Close()

	type orderResp struct {
		ID              int64           `json:"id"`
		UserID          int64           `json:"user_id"`
		RestaurantID    int64           `json:"restaurant_id"`
		Items           json.RawMessage `json:"items"`
		Total           float64         `json:"total"`
		Status          string          `json:"status"`
		Address         string          `json:"address"`
		Note            string          `json:"note"`
		CreatedAt       string          `json:"created_at"`
		UpdatedAt       string          `json:"updated_at"`
		RestaurantName  string          `json:"restaurant_name"`
		RestaurantImage string          `json:"restaurant_image"`
	}

	var orders []orderResp
	for rows.Next() {
		var o orderResp
		var items []byte
		if err := rows.Scan(&o.ID, &o.UserID, &o.RestaurantID, &items, &o.Total,
			&o.Status, &o.Address, &o.Note, &o.CreatedAt, &o.UpdatedAt,
			&o.RestaurantName, &o.RestaurantImage); err != nil {
			mw.Error(w, http.StatusInternalServerError, "扫描订单数据失败")
			return
		}
		if items == nil {
			o.Items = json.RawMessage("[]")
		} else {
			o.Items = json.RawMessage(items)
		}
		orders = append(orders, o)
	}
	if orders == nil {
		orders = []orderResp{}
	}

	mw.JSON(w, http.StatusOK, map[string]interface{}{"data": orders})
}

// POST /api/orders
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())

	var body struct {
		RestaurantID int64           `json:"restaurant_id"`
		Items        json.RawMessage `json:"items"`
		Total        float64         `json:"total"`
		Address      string          `json:"address"`
		Note         string          `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "无效的请求体")
		return
	}
	if body.RestaurantID == 0 {
		mw.Error(w, http.StatusBadRequest, "请选择餐厅")
		return
	}
	if body.Items == nil {
		body.Items = json.RawMessage("[]")
	}

	// 验证餐厅存在
	var exists bool
	_ = h.DB.Pool.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM restaurants WHERE id=$1)`, body.RestaurantID).Scan(&exists)
	if !exists {
		mw.Error(w, http.StatusBadRequest, "餐厅不存在")
		return
	}

	var id int64
	err := h.DB.Pool.QueryRow(r.Context(), `
		INSERT INTO orders (user_id, restaurant_id, items, total, status, address, note)
		VALUES ($1, $2, $3, $4, 'pending', $5, $6)
		RETURNING id
	`, userID, body.RestaurantID, []byte(body.Items), body.Total, body.Address, body.Note).Scan(&id)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "创建订单失败")
		return
	}

	mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "订单创建成功"})
}

// GET /api/orders/{id}
func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var o struct {
		ID              int64           `json:"id"`
		UserID          int64           `json:"user_id"`
		RestaurantID    int64           `json:"restaurant_id"`
		Items           json.RawMessage `json:"items"`
		Total           float64         `json:"total"`
		Status          string          `json:"status"`
		Address         string          `json:"address"`
		Note            string          `json:"note"`
		CreatedAt       string          `json:"created_at"`
		UpdatedAt       string          `json:"updated_at"`
		RestaurantName  string          `json:"restaurant_name"`
		RestaurantImage string          `json:"restaurant_image"`
	}

	var items []byte
	err = h.DB.Pool.QueryRow(r.Context(), `
		SELECT o.id, o.user_id, o.restaurant_id, o.items, o.total, o.status, o.address, o.note,
		       o.created_at, o.updated_at,
		       r.name, r.image_url
		FROM orders o
		LEFT JOIN restaurants r ON r.id = o.restaurant_id
		WHERE o.id = $1 AND o.user_id = $2
	`, id, userID).Scan(&o.ID, &o.UserID, &o.RestaurantID, &items, &o.Total,
		&o.Status, &o.Address, &o.Note, &o.CreatedAt, &o.UpdatedAt,
		&o.RestaurantName, &o.RestaurantImage)
	if err == pgx.ErrNoRows {
		mw.Error(w, http.StatusNotFound, "订单不存在")
		return
	}
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "查询订单失败")
		return
	}
	if items == nil {
		o.Items = json.RawMessage("[]")
	} else {
		o.Items = json.RawMessage(items)
	}

	mw.JSON(w, http.StatusOK, o)
}

// PATCH /api/orders/{id}/status
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := mw.GetUserID(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		mw.Error(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		mw.Error(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	// 允许的状态: pending, confirmed, preparing, delivering, delivered, cancelled
	validStatuses := map[string]bool{
		"pending":    true,
		"confirmed":  true,
		"preparing":  true,
		"delivering": true,
		"delivered":  true,
		"cancelled":  true,
	}
	if !validStatuses[body.Status] {
		mw.Error(w, http.StatusBadRequest, "无效的订单状态")
		return
	}

	result, err := h.DB.Pool.Exec(r.Context(), `
		UPDATE orders SET status = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`, body.Status, id, userID)
	if err != nil {
		mw.Error(w, http.StatusInternalServerError, "更新订单状态失败")
		return
	}
	if result.RowsAffected() == 0 {
		mw.Error(w, http.StatusNotFound, "订单不存在")
		return
	}

	mw.JSON(w, http.StatusOK, map[string]string{"message": "订单状态已更新"})
}
