package main

import (
	"fmt"
	"github.com/koloo91/monhttp/controller"
	"github.com/koloo91/monhttp/notifier"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting monhttp")

	log.Info("Loading configuration")
	if err := service.LoadConfig(); err != nil {
		log.Fatalf("Unable to load configuration: '%s'", err)
	}

	notificationSystem := notifier.NewNotificationSystem()
	notificationSystem.SetupDefaultNotifier()
	notificationSystem.Start()

	service.SetNotificationSystem(notificationSystem)

	service.LoadUsers()

	if service.IsSetup() {
		host := service.GetConfig().Host
		port := service.GetConfig().Port
		user := service.GetConfig().User
		password := service.GetConfig().Password
		databaseName := service.GetConfig().DatabaseName

		log.Info("Service is setup. Connecting to database")
		err := service.LoadDatabase(host, port, user, password, databaseName, "./migrations")
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
		Addr:         fmt.Sprintf(":%d", service.GetConfig().ServerPort),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      router,
	}

	log.Infof("Starting http server on port '%d'", service.GetConfig().ServerPort)
	log.Fatal(server.ListenAndServe())
}
