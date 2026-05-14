package handlers

import (
	"encoding/json"
	"net/http"
	"todo-api/internal/auth"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := auth.HashPassword(user.Password)

	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?) RETURNING id`
	err := database.DB.QueryRow(query, user.Name, user.Email, hashedPassword).Scan(&user.ID)

	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}
	token, _ := auth.GenerateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	var dbUser models.User
	json.NewDecoder(r.Body).Decode(&input)

	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", input.Email).Scan(&dbUser.ID, &dbUser.Password)
	if err != nil || !auth.CheckPasswordHash(input.Password, dbUser.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(dbUser.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
