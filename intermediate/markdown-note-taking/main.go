package main

import (
	"log"
	"net/http"

	"markdown-note-taking/handler"
	"markdown-note-taking/notes"
)

func main() {
	repo, err := notes.NewFileRepository("./storage")
	if err != nil {
		log.Fatalf("Failed to initialize system file repository layer: %v", err)
	}
	noteService := notes.NewService(repo)
	httpHandler := handler.NewHTTPHandler(noteService)

	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	log.Println("Launching Markdown Note-taking application on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server connection loop aborted: %v", err)
	}
}
