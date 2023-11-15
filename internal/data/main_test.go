//go:build integration
// +build integration

package data

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testModels Models

func TestMain(m *testing.M) {
	var dsn string
	var err error

	flag.StringVar(&dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.Parse()

	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln("cannot connect to database: ", err)
	}

	testModels = NewModels(testDB)

	os.Exit(m.Run())
}
