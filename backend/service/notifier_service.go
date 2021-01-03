package service

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetNotifiers() []model.Notify {
	return notificationSystem.GetNotifiers()
}

func UpdateNotifier(id string, body map[string]interface{}) error {
	log.Infof("Updating notififier with id '%s'", id)
	for k, v := range body {
		viper.Set(fmt.Sprintf("notifier.%s.%s", id, k), v)
	}

	notificationSystem.SetupDefaultNotifier()

	return viper.WriteConfig()
}
