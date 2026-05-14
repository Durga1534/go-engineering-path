package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("sqlite", "./todo.db")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign keys: ", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT UNIQUE,
	password TEXT
);	

    CREATE TABLE IF NOT EXISTS todos (
        id  INTEGER PRIMARY KEY AUTOINCREMENT,
		user_Id INTEGER NOT NULL,
		title  TEXT NOT NULL,
		description  TEXT,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	
	);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables: ", err)
	}
	log.Println("Database initialized successfully with Foreign Keys enabled.")
}
