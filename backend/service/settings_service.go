package service

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	isSetup bool
)

func IsSetup() bool {
	// TODO: refactor this
	return isSetup
}

func SetIsSetup(b bool) {
	isSetup = b
}

func UpdateSettings(settings model.SettingsVo) error {
	viper.Set("database.host", settings.DatabaseHost)
	viper.Set("database.port", settings.DatabasePort)
	viper.Set("database.user", settings.DatabaseUser)
	viper.Set("database.password", settings.DatabasePassword)
	viper.Set("database.name", settings.DatabaseName)

	viper.Set("users", []string{fmt.Sprintf("%s:%s", settings.Username, settings.Password)})

	host := settings.DatabaseHost
	port := settings.DatabasePort
	user := settings.DatabaseUser
	password := settings.DatabasePassword
	databaseName := settings.DatabaseName

	if err := LoadDatabase(host, port, user, password, databaseName); err != nil {
		log.Errorf("Unable to load database with configuration: '%s'", err)
		return err
	}

	if err := AddUser(settings.Username, settings.Password); err != nil {
		return err
	}

	log.Info("Writing new settings into config.yml")
	return viper.WriteConfigAs("./config/config.yml")
}
