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

func connectToDatabase() (*sql.DB, error) {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.name")

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

func LoadDatabase() error {
	if database != nil {
		return nil
	}

	var err error
	database, err = connectToDatabase()
	if err != nil {
		log.Errorf("Unable to connect to database: '%s'", err)
		return err
	}

	if err := runDatabaseMigrations(database); err != nil {
		log.Errorf("Unable to run database migrations: '%s'", err)
		return err
	}

	repository.SetDatabase(database)
	go StartScheduleJob()

	SetIsSetup(true)
	return nil
}

func GetDatabase() *sql.DB {
	return database
}
