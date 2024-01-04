package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	connStr := "user=bank-transfer dbname=postgres password=123456 port=5432 host=localhost search_path=bank sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type Queries struct {
	DB *sql.DB
}
