package service

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/koloo91/monhttp/repository"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	database *sql.DB
)

func connectToDatabase(host string, port int, user, password, databaseName string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, databaseName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func runDatabaseMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		log.Warn(err)
	}
	return nil
}

func LoadDatabase(host string, port int, user, password, databaseName string) error {
	if database != nil {
		return nil
	}

	var err error
	database, err = connectToDatabase(host, port, user, password, databaseName)
	if err != nil {
		log.Errorf("Unable to connect to database: '%s'", err)
		return err
	}

	if err := runDatabaseMigrations(database); err != nil {
		log.Errorf("Unable to run database migrations: '%s'", err)
		return err
	}

	repository.SetDatabase(database)
	go StartScheduleJob(viper.GetBool("scheduler.enabled"))

	return nil
}

func GetDatabase() *sql.DB {
	return database
}
