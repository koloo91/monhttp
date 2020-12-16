package repository

import "database/sql"

var db *sql.DB

func SetDatabase(database *sql.DB) {
	db = database

	prepareServiceStatements()
}

func BeginnTransaction() (*sql.Tx, error) {
	return db.Begin()
}
