package db

import (
	"database/sql"
	"log"
	"test-task/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() *sql.DB {
	dsn := config.DBDSN()

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
