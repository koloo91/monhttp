package service

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	"github.com/spf13/viper"
)

func GetNotifiers() []model.Notify {
	return notificationSystem.GetNotifiers()
}

func UpdateNotifier(id string, body map[string]interface{}) error {
	for k, v := range body {
		viper.Set(fmt.Sprintf("notifier.%s.%s", id, k), v)
	}
	return viper.WriteConfig()
}
