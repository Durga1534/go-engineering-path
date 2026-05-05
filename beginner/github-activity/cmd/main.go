package main

import (
	"fmt"
	"github-activity/internal/github"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	username := os.Args[1]
	events, err := github.FetchActivity(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(events) == 0 {
		fmt.Println("No recent activity found for this user.")
		return
	}

	fmt.Printf("Recent activity for %s:\n", username)
	for _, event := range events {
		displayEvent(event)
	}
}

func displayEvent(event github.GitHubEvent) {
	switch event.Type {
	case "PushEvent":
		fmt.Printf("- Pushed %d commits to %s\n", len(event.Payload.Commits), event.Repo.Name)
	case "IssuesEvent":
		fmt.Printf("- Opened/Closed an issue in %s\n", event.Repo.Name)
	case "WatchEvent":
		fmt.Printf("- Starred %s\n", event.Repo.Name)
	case "CreateEvent":
		fmt.Printf("- Created a new resource in %s\n", event.Repo.Name)
	default:
		fmt.Printf("- %s in %s\n", event.Type, event.Repo.Name)
	}
}
