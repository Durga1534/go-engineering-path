package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	DB, err = sql.Open("sqlite", "./expense.db")
	if err != nil {
		log.Fatal("DB Connection failed: ", err)
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Foreign keys execution failed: ", err)
	}
	query := `
	CREATE TABLE IF NOT EXISTS users (
	        id    INTEGER PRIMARY KEY AUTOINCREMENT,
			name  TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
	);   
	CREATE TABLE IF NOT EXISTS expenses (
	      id   INTEGER PRIMARY KEY AUTOINCREMENT,
		  user_id INTEGER NOT NULL,
		  title  TEXT NOT NULL,
		  amount  REAL NOT NULL,
		  category TEXT NOT NULL CHECK(category IN ('Groceries', 'Leisure', 'Electronics', 'Utilities', 'Clothing', 'Health', 'Others')),
		  date  DATETIME NOT NULL,
		  created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal("Table compilation failed:", err)
	}
	log.Println("Expense Tracker Database Ready")
}
