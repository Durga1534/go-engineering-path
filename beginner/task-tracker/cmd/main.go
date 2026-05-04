package main

import (
	"fmt"
	"os"
	"strconv"
	"task-tracker/internal/task"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli [command] [arguments]")
		return
	}

	command := os.Args[1]
	tasks, _ := task.LoadTasks()

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task description.")
			return
		}
		newID := 1
		if len(tasks) > 0 {
			newID = tasks[len(tasks)-1].ID + 1
		}

		newTask := task.Task{
			ID:          newID,
			Description: os.Args[2],
			Status:      "todo",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		tasks = append(tasks, newTask)
		task.SaveTasks(tasks)
		fmt.Printf("Task added successfully (ID: %d)\n", newID)

	case "list":
		for _, t := range tasks {
			fmt.Printf("[%d] %s - %s\n", t.ID, t.Description, t.Status)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide an ID")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		var updatedTasks []task.Task
		for _, t := range tasks {
			if t.ID != id {
				updatedTasks = append(updatedTasks, t)
			}
		}
		task.SaveTasks(updatedTasks)
		fmt.Println("Task deleted.")

	default:
		fmt.Println("Unknown command.")
	}
}
