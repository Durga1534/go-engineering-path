package main

import (
	"blogging-api/database"
	"blogging-api/internal/handlers"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	database.ConnectDB()

	defer database.DB.Close()
	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPosts(w, r)
		case http.MethodPost:
			handlers.CreatePost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/posts/")
		if id == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handlers.GetPostByID(w, r, id)
		case http.MethodPut:
			handlers.UpdatePost(w, r, id)
		case http.MethodDelete:
			handlers.DeletePost(w, r, id)
		}
	})

	fmt.Println("REST API is running at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
