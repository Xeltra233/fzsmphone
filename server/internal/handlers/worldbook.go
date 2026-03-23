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

type WorldBookHandler struct {
    DB *database.DB
}

type worldBookRecord struct {
    ID        int64           `json:"id"`
    UserID    int64           `json:"user_id"`
    Name      string          `json:"name"`
    BindChars json.RawMessage `json:"bind_chars"`
    Entries   json.RawMessage `json:"entries"`
    CreatedAt string          `json:"created_at"`
    UpdatedAt string          `json:"updated_at"`
}

// GET /api/worldbook
func (h *WorldBookHandler) List(w http.ResponseWriter, r *http.Request) {
    userID, _ := mw.GetUserID(r.Context())

    rows, err := h.DB.Pool.Query(r.Context(), `
        SELECT id, user_id, name, bind_chars, entries, created_at, updated_at
        FROM worldbook_books WHERE user_id = $1
        ORDER BY created_at DESC
    `, userID)
    if err != nil {
        mw.Error(w, http.StatusInternalServerError, "failed to query worldbooks")
        return
    }
    defer rows.Close()

    books := make([]worldBookRecord, 0)
    for rows.Next() {
        var book worldBookRecord
        if err := rows.Scan(&book.ID, &book.UserID, &book.Name, &book.BindChars, &book.Entries, &book.CreatedAt, &book.UpdatedAt); err != nil {
            mw.Error(w, http.StatusInternalServerError, "failed to scan worldbook")
            return
        }
        if len(book.BindChars) == 0 {
            book.BindChars = json.RawMessage("[]")
        }
        if len(book.Entries) == 0 {
            book.Entries = json.RawMessage("[]")
        }
        books = append(books, book)
    }

    mw.JSON(w, http.StatusOK, map[string]interface{}{"data": books})
}

// POST /api/worldbook
func (h *WorldBookHandler) Create(w http.ResponseWriter, r *http.Request) {
    userID, _ := mw.GetUserID(r.Context())

    var body struct {
        Name      string          `json:"name"`
        BindChars json.RawMessage `json:"bind_chars"`
        Entries   json.RawMessage `json:"entries"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        mw.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }
    if body.Name == "" {
        mw.Error(w, http.StatusBadRequest, "name is required")
        return
    }
    if len(body.BindChars) == 0 {
        body.BindChars = json.RawMessage("[]")
    }
    if len(body.Entries) == 0 {
        body.Entries = json.RawMessage("[]")
    }

    var id int64
    err := h.DB.Pool.QueryRow(r.Context(), `
        INSERT INTO worldbook_books (user_id, name, bind_chars, entries)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, userID, body.Name, body.BindChars, body.Entries).Scan(&id)
    if err != nil {
        mw.Error(w, http.StatusInternalServerError, "failed to create worldbook")
        return
    }

    mw.JSON(w, http.StatusCreated, map[string]interface{}{"id": id, "message": "worldbook created"})
}

// GET /api/worldbook/{id}
func (h *WorldBookHandler) Get(w http.ResponseWriter, r *http.Request) {
    userID, _ := mw.GetUserID(r.Context())
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        mw.Error(w, http.StatusBadRequest, "invalid id")
        return
    }

    var book worldBookRecord
    err = h.DB.Pool.QueryRow(r.Context(), `
        SELECT id, user_id, name, bind_chars, entries, created_at, updated_at
        FROM worldbook_books WHERE id = $1 AND user_id = $2
    `, id, userID).Scan(&book.ID, &book.UserID, &book.Name, &book.BindChars, &book.Entries, &book.CreatedAt, &book.UpdatedAt)
    if err == pgx.ErrNoRows {
        mw.Error(w, http.StatusNotFound, "worldbook not found")
        return
    }
    if err != nil {
        mw.Error(w, http.StatusInternalServerError, "failed to get worldbook")
        return
    }
    if len(book.BindChars) == 0 {
        book.BindChars = json.RawMessage("[]")
    }
    if len(book.Entries) == 0 {
        book.Entries = json.RawMessage("[]")
    }

    mw.JSON(w, http.StatusOK, book)
}

// PUT /api/worldbook/{id}
func (h *WorldBookHandler) Update(w http.ResponseWriter, r *http.Request) {
    userID, _ := mw.GetUserID(r.Context())
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        mw.Error(w, http.StatusBadRequest, "invalid id")
        return
    }

    var body struct {
        Name      string          `json:"name"`
        BindChars json.RawMessage `json:"bind_chars"`
        Entries   json.RawMessage `json:"entries"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        mw.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }
    if len(body.BindChars) == 0 {
        body.BindChars = json.RawMessage("[]")
    }
    if len(body.Entries) == 0 {
        body.Entries = json.RawMessage("[]")
    }

    result, err := h.DB.Pool.Exec(r.Context(), `
        UPDATE worldbook_books SET name=$1, bind_chars=$2, entries=$3, updated_at=NOW()
        WHERE id=$4 AND user_id=$5
    `, body.Name, body.BindChars, body.Entries, id, userID)
    if err != nil {
        mw.Error(w, http.StatusInternalServerError, "failed to update worldbook")
        return
    }
    if result.RowsAffected() == 0 {
        mw.Error(w, http.StatusNotFound, "worldbook not found")
        return
    }

    mw.JSON(w, http.StatusOK, map[string]string{"message": "worldbook updated"})
}

// DELETE /api/worldbook/{id}
func (h *WorldBookHandler) Delete(w http.ResponseWriter, r *http.Request) {
    userID, _ := mw.GetUserID(r.Context())
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        mw.Error(w, http.StatusBadRequest, "invalid id")
        return
    }

    result, err := h.DB.Pool.Exec(r.Context(), `
        DELETE FROM worldbook_books WHERE id = $1 AND user_id = $2
    `, id, userID)
    if err != nil {
        mw.Error(w, http.StatusInternalServerError, "failed to delete worldbook")
        return
    }
    if result.RowsAffected() == 0 {
        mw.Error(w, http.StatusNotFound, "worldbook not found")
        return
    }

    mw.JSON(w, http.StatusOK, map[string]string{"message": "worldbook deleted"})
}
