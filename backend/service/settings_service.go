package service

import (
	"github.com/koloo91/monhttp/model"
	"github.com/spf13/viper"
)

var (
	isSetup bool
)

func IsSetup() bool {
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

	if err := LoadDatabase(); err != nil {
		return err
	}

	return viper.WriteConfigAs("config.yml")
}
