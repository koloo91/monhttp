package service

import (
	"github.com/koloo91/monhttp/repository"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func AddUser(name, password string) error {
	return repository.InsertUser(name, password)
}

func GetUsers() map[string]string {
	return repository.GetUsers()
}

func RemoveUser(name string) error {
	return repository.RemoveUser(name)
}

func LoadUsersFromConfig() {

	users := viper.GetStringSlice("users")
	for _, user := range users {
		usernameAndPassword := strings.Split(user, ":")
		if len(usernameAndPassword) != 2 {
			continue
		}

		if err := AddUser(usernameAndPassword[0], usernameAndPassword[1]); err != nil {
			log.Errorf("Unable to add user: '%s'", err)
		}
	}
}
