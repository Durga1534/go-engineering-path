package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"todo-api/internal/database"
	"todo-api/internal/handlers"
)

func main() {
	// 1. Initialize Database
	database.ConnectDB()
	defer database.DB.Close()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	http.HandleFunc("/todos", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTodos(w, r)
		case http.MethodPost:
			handlers.CreateTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/todos/", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
		if idStr == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodPut:
			handlers.UpdateTodo(w, r, idStr)
		case http.MethodDelete:
			handlers.DeleteTodo(w, r, idStr)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	port := ":8000"
	fmt.Printf(" Todo API is running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
