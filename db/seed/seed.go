package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	var dsn string

	flag.StringVar(&dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.Parse()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	script, err := os.ReadFile("./db/seed/movies.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seed completed 🌱")
}
