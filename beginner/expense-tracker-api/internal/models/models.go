package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Expense struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}
