package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:HJFHSN2OFJFEC52UKV2CMNBBKQ@db:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	for db.Ping() != nil {
		log.Println("Waiting for database...")
		time.Sleep(time.Second)
	}
	log.Println("Database is available")

	err = MigrateDatabase(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database successfully migrated")

	log.Println("Listening for HTTP connections...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Shutting down")
}
