package db

import (
	"database/sql"
	"log"
	"os"
)

func InitDB() *sql.DB {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Println("postgres db: Use default settings")
		dsn = "postgres://myadmin:mypassword@localhost:5432/mydatabase"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Error: failed to prepare database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error: faild to connect with db: %v", err)
	}

	log.Println("Connected to Data Base")
	return db
}
