package service

import (
	"github.com/koloo91/monhttp/model"
)

func GetNotifiers() []model.Notify {
	return notificationSystem.GetNotifiers()
}
