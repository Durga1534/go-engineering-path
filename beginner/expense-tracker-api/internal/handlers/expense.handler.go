package handlers

import (
	"database/sql"
	"encoding/json"
	"expense-tracker-api/internal/database"
	"expense-tracker-api/internal/models"
	"net/http"
	"strconv"
	"time"
)

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	var exp models.Expense
	json.NewDecoder(r.Body).Decode(&exp)

	if exp.Date.IsZero() {
		exp.Date = time.Now()
	}

	query := `INSERT INTO expenses (user_id, title, amount, category, date) VALUES(?,?,?,?) RETURNING id`
	err := database.DB.QueryRow(query, userID, exp.Title, exp.Amount, exp.Category, exp.Date).Scan(&exp.ID)
	if err != nil {
		http.Error(w, "Invalid Category or Database constraint violation", http.StatusBadRequest)
		return
	}

	exp.UserID = userID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exp)
}

func GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	filter := r.URL.Query().Get("filter")

	baseQuery := `SELECT id, title, amount, category, date FROM expenses WHERE user_id = ?`
	var args []interface{}
	args = append(args, userID)

	switch filter {
	case "week":
		baseQuery += ` AND date >= datetime('now', '-7 days')`
	case "month":
		baseQuery += ` AND date >= datetime('now', '-1 month')`
	case "three_months":
		baseQuery += ` AND date >= datetime('now', '-3 months')`
	case "custom":
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		baseQuery += ` AND date BETWEEN ? AND ?`
		args = append(args, start+" 00:00:00", end+" 23:59:59")
	}

	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		var exp models.Expense
		var dateStr string
		rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Category, &dateStr)
		exp.Date, _ = time.Parse("2006-01-02T15:04:05Z", dateStr)
		expenses = append(expenses, exp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func UpdateExpense(w http.ResponseWriter, r *http.Request, idStr string) {
	userID := r.Context().Value("user_id").(int)
	var exp models.Expense
	json.NewDecoder(r.Body).Decode(&exp)

	var dbUserID int
	err := database.DB.QueryRow("SELECT user_id FROM expenses WHERE id = ?", idStr).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if dbUserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	query := `UPDATE expenses SET title = ?, amount = ?, category = ?, date = ? WHERE id = ?`
	_, err = database.DB.Exec(query, exp.Title, exp.Amount, exp.Category, exp.Date, idStr)
	if err != nil {
		http.Error(w, "Update rejected due to validation failure", http.StatusBadRequest)
		return
	}
	exp.ID, _ = strconv.Atoi(idStr)
	json.NewEncoder(w).Encode(exp)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request, idStr string) {
	userID := r.Context().Value("user_id").(int)

	var dbUserID int
	err := database.DB.QueryRow("SELECT user_id FROM expenses WHERE id = ? ", idStr).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if dbUserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	database.DB.Exec("DELETE FROM expenses WHERE id = ?", idStr)
	w.WriteHeader(http.StatusNoContent)
}
