package handlers

import (
	"encoding/json"
	"expense-tracker-api/internal/auth"
	"expense-tracker-api/internal/database"
	"expense-tracker-api/internal/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := auth.HashPassword(user.Password)

	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?) RETURNING id`
	err := database.DB.QueryRow(query, user.Name, user.Email, hashedPassword).Scan(&user.ID)

	if err != nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	token, _ := auth.GenerateToken(user.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.User
	var dbUser models.User
	json.NewDecoder(r.Body).Decode(&input)

	query := `SELECT id, password FROM users WHERE email = ?`
	err := database.DB.QueryRow(query, input.Email).Scan(&dbUser.ID, &dbUser.Password)
	if err != nil || !auth.CheckPasswordHash(input.Password, dbUser.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(dbUser.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
