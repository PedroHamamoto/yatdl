package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // postgres driver, used by database/sql
)

func ConnectToDatabase() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
