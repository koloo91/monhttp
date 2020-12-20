package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/koloo91/monhttp/controller"
	"github.com/koloo91/monhttp/notifier"
	"github.com/koloo91/monhttp/service"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	service.SetIsSetup(true)

	if err := viper.ReadInConfig(); err != nil {
		service.SetIsSetup(false)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("No config file found")
		} else {
			log.Fatal(err)
		}
	}

	notificationSystem := notifier.NewNotificationSystem()
	notificationSystem.SetupDefaultNotifier()
	notificationSystem.Start()

	service.SetNotificationSystem(notificationSystem)

	if service.IsSetup() {
		err := service.LoadDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}

	defer func() {
		if service.GetDatabase() != nil {
			service.GetDatabase().Close()
		}
	}()

	router := controller.SetupRoutes()

	server := http.Server{
		Addr:         ":8081",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      router,
	}

	log.Fatal(server.ListenAndServe())
}
