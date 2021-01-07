package service

import (
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config model.Config
)

func GetConfig() model.Config {
	return config
}

func LoadConfig() error {
	viper.AddConfigPath("./config")
	viper.WatchConfig()
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.SetDefault("SERVER_PORT", 8081)
	viper.SetDefault("SCHEDULER_NUMBER_OF_WORKERS", 5)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("No config file found")
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}
