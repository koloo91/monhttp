package service

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var (
	config model.Config
)

func GetConfig() model.Config {
	return config
}

func LoadConfig() error {
	configLocation := "./config"
	configName := "config"

	_ = os.Mkdir(configLocation, os.ModePerm)

	viper.AddConfigPath(configLocation)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.SetDefault("DATABASE_HOST", "")
	viper.SetDefault("DATABASE_PORT", 0)
	viper.SetDefault("DATABASE_USER", "")
	viper.SetDefault("DATABASE_PASSWORD", "")
	viper.SetDefault("DATABASE_NAME", "")

	viper.SetDefault("USERS", "")

	viper.SetDefault("SERVER_PORT", 8081)
	viper.SetDefault("SCHEDULER_ENABLED", true)
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

	return viper.WriteConfigAs(fmt.Sprintf("%s/%s.env", configLocation, configName))
}
