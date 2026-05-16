package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"expense-tracker-api/internal/database"
	"expense-tracker-api/internal/handlers"
	"expense-tracker-api/internal/middlewares"
)

func main() {
	database.ConnectDB()
	defer database.DB.Close()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	http.HandleFunc("/expenses", middlewares.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetExpenses(w, r)
		case http.MethodPost:
			handlers.CreateExpense(w, r)
		default:
			http.Error(w, "Method configuration unallowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/expenses/", middlewares.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
		if idStr == "" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodPut:
			handlers.UpdateExpense(w, r, idStr)
		case http.MethodDelete:
			handlers.DeleteExpense(w, r, idStr)
		default:
			http.Error(w, "Method configuartion unallowed", http.StatusMethodNotAllowed)
		}
	}))

	port := ":5000"
	fmt.Printf("Final Beginner Challenge Running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
