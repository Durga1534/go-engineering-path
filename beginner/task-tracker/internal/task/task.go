package task

import (
	"encoding/json"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const fileName = "tasks.json"

func LoadTasks() ([]Task, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}
