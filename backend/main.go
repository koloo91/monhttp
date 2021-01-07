package main

import (
	"fmt"
	"github.com/koloo91/monhttp/controller"
	"github.com/koloo91/monhttp/notifier"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting monhttp")

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	viper.WatchConfig()

	viper.SetDefault("server.port", 8081)
	viper.SetDefault("scheduler.enabled", true)
	viper.SetDefault("scheduler.numberOfWorkers", 5)

	service.SetIsSetup(true)

	log.Info("Reading configuration file")
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
		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		user := viper.GetString("database.user")
		password := viper.GetString("database.password")
		databaseName := viper.GetString("database.name")

		err := service.LoadDatabase(host, port, user, password, databaseName)
		if err != nil {
			log.Fatalf("Unable to connect to database: '%s'", err)
		}
	}

	defer func() {
		if service.GetDatabase() != nil {
			service.GetDatabase().Close()
		}
	}()

	router := controller.SetupRoutes()

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", viper.GetInt("server.port")),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      router,
	}

	log.Infof("Starting http server on port '%d'", viper.GetInt("server.port"))
	log.Fatal(server.ListenAndServe())
}
