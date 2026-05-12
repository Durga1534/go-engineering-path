package handlers

import (
	"blogging-api/database"
	"blogging-api/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// CreatePost - POST /posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var p models.Post

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if p.Title == "" || p.Content == "" || p.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title, Content, and Category are required"})
		return
	}

	now := time.Now()
	tagsString := strings.Join(p.Tags, ",")

	query := `INSERT INTO posts (title, content, category, tags, createdAt, updatedAt) 
	          VALUES (?, ?, ?, ?, ?, ?) RETURNING id`

	err := database.DB.QueryRow(query, p.Title, p.Content, p.Category, tagsString, now, now).Scan(&p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.CreatedAt = now
	p.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// GetPosts - GET /posts (includes search filter)
func GetPosts(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	var rows *sql.Rows
	var err error

	if term != "" {
		query := "SELECT id, title, content, category, tags, createdAt, updatedAt FROM posts WHERE title LIKE ? OR content LIKE ? OR category LIKE ?"
		wildcard := "%" + term + "%"
		rows, err = database.DB.Query(query, wildcard, wildcard, wildcard)
	} else {
		rows, err = database.DB.Query("SELECT id, title, content, category, tags, createdAt, updatedAt FROM posts")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		var tagsStr string
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &tagsStr, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			continue
		}
		p.Tags = strings.Split(tagsStr, ",")
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetPostByID - GET /posts/{id}
func GetPostByID(w http.ResponseWriter, r *http.Request, id string) {
	var p models.Post
	var tagsStr string

	query := "SELECT id, title, content, category, tags, createdAt, updatedAt FROM posts WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &p.Category, &tagsStr, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.Tags = strings.Split(tagsStr, ",")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdatePost - PUT /posts/{id}
func UpdatePost(w http.ResponseWriter, r *http.Request, id string) {
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	tagsStr := strings.Join(p.Tags, ",")

	query := `UPDATE posts SET title = ?, content = ?, category = ?, tags = ?, updatedAt = ? WHERE id = ?`
	result, err := database.DB.Exec(query, p.Title, p.Content, p.Category, tagsStr, now, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	// Return updated post (fetch it back or construct it)
	p.UpdatedAt = now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DeletePost - DELETE /posts/{id}
func DeletePost(w http.ResponseWriter, r *http.Request, id string) {
	query := "DELETE FROM posts WHERE id = ?"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
