package repository

import "database/sql"

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