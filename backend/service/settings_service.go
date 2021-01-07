package service

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	onAdminSetCallback func(string, string)
)

func IsSetup() bool {
	return len(config.Host) > 0 &&
		config.Port > 0 &&
		len(config.User) > 0 &&
		len(config.Password) > 0 &&
		len(config.DatabaseName) > 0 &&
		len(config.Users) > 0
}

func UpdateSettings(settings model.SettingsVo) error {
	viper.Set("DATABASE_HOST", settings.DatabaseHost)
	viper.Set("DATABASE_PORT", settings.DatabasePort)
	viper.Set("DATABASE_USER", settings.DatabaseUser)
	viper.Set("DATABASE_PASSWORD", settings.DatabasePassword)
	viper.Set("DATABASE_NAME", settings.DatabaseName)

	viper.Set("USERS", fmt.Sprintf("%s:%s", settings.Username, settings.Password))

	if err := LoadConfig(); err != nil {
		return err
	}

	if err := LoadDatabase(); err != nil {
		log.Errorf("Unable to load database with configuration: '%s'", err)
		return err
	}

	onAdminSetCallback(settings.Username, settings.Password)

	log.Info("Writing new settings into config.env")
	return viper.WriteConfigAs("./config/config.env")
}

func LoadUsers() map[string]string {
	usersMap := make(map[string]string)

	users := strings.Split(GetConfig().Users, ",")
	for _, user := range users {
		usernameAndPassword := strings.Split(user, ":")
		if len(usernameAndPassword) != 2 {
			continue
		}
		usersMap[usernameAndPassword[0]] = usernameAndPassword[1]
	}

	return usersMap
}
func SetOnAdminSetCallback(callback func(string, string)) {
	onAdminSetCallback = callback
}
