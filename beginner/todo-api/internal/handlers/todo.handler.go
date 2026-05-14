package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO todos (user_id, title, description, created_at) VALUES (?,?,?,?) RETURNING id`
	now := time.Now()
	err := database.DB.QueryRow(query, userID, todo.Title, todo.Description, now).Scan(&todo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.UserID = userID
	todo.CreatedAt = now
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)

}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, _ := database.DB.Query("SELECT id, title, description FROM todos WHERE user_id = ? LIMIT ? OFFSET ?", userID, limit, offset)

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		rows.Scan(&t.ID, &t.Title, &t.Description)
		todos = append(todos, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request, idStr string) {
	userID := r.Context().Value("user_id").(int)
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	var dbUserID int
	err := database.DB.QueryRow("SELECT user_id FROMM todos WHERE id=?", idStr).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if dbUserID != userID {
		http.Error(w, "Forbidden: You do not own this", http.StatusForbidden)
		return
	}
	query := `UPDATE todos SET title=?, description=? WHERE id = ?`
	_, err = database.DB.Exec(query, todo.Title, todo.Description, idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.ID, _ = strconv.Atoi(idStr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request, idStr string) {
	userID := r.Context().Value("user_id").(int)

	var dbUserID int
	err := database.DB.QueryRow("SELECT user_id FROM todos WHERE id = ?", idStr).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if dbUserID != userID {
		http.Error(w, "Forbidden: ", http.StatusForbidden)
		return
	}
	_, err = database.DB.Exec("DELETE FROM todos WHERE id = ?", idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
