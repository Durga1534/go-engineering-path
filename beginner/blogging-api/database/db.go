package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("sqlite", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}

	query := `
    CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    content TEXT,
    category TEXT,
    tags TEXT,
    createdAt DATETIME,
    updatedAt DATETIME
);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	alterQuery := `ALTER TABLE posts ADD COLUMN updatedAt DATETIME;`
	DB.Exec(alterQuery)
}
