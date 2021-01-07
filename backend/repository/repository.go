package repository

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var db *sql.DB

func SetDatabase(database *sql.DB) {
	db = database

	prepareServiceStatements()
	prepareCheckStatements()
	prepareFailureStatements()
}

func BeginnTransaction() (*sql.Tx, error) {
	return db.Begin()
}
