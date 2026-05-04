# Task Tracker CLI

A simple Command Line Interface (CLI) built with **Go** to manage your daily tasks. This project is part of my journey to master backend engineering and follows the roadmap and requirements from [roadmap.sh](https://roadmap.sh/projects/task-tracker).

## 🚀 Features
- **Add Tasks**: Quickly create new tasks with descriptions.
- **List Tasks**: View all your tasks, their status, and unique IDs.
- **Delete Tasks**: Remove tasks from your list using their ID.
- **Persistence**: All data is saved to a local `tasks.json` file.
- **No External Dependencies**: Built using only Go's standard library to master the fundamentals.

## 🛠️ Tech Stack & Learning Goals
- **Language**: Go (Golang)
- **Concepts Learned**:
  - Structuring Go projects using the `/cmd` and `/internal` pattern.
  - Working with the File System (`os` package).
  - JSON encoding/decoding (`encoding/json`).
  - Handling CLI arguments and user input.

## 📁 Project Structure
```text
task-tracker/
├── cmd/
│   └── main.go       # Entry point & CLI Routing
├── internal/
│   └── task/
│       ├── task.go   # Data structures & Blueprint
│       └── store.go  # File I/O & Business Logic
└── tasks.json        # Local database (Auto-generated)