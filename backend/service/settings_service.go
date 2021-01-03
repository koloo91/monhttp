package service

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	isSetup            bool
	onAdminSetCallback func(string, string)
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

	viper.Set("users", []string{fmt.Sprintf("%s:%s", settings.Username, settings.Password)})

	if err := LoadDatabase(); err != nil {
		log.Errorf("Unable to load database with configuration: '%s'", err)
		return err
	}

	onAdminSetCallback(settings.Username, settings.Password)

	log.Info("Writing new settings into config.yml")
	return viper.WriteConfigAs("./config/config.yml")
}

func LoadUsers() map[string]string {
	usersMap := make(map[string]string)

	users := viper.GetStringSlice("users")
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
