package model

type Config struct {
	ServerPort int `mapstructure:"SERVER_PORT"`

	SchedulerEnabled         bool `mapstructure:"SCHEDULER_ENABLED"`
	SchedulerNumberOfWorkers bool `mapstructure:"SCHEDULER_NUMBER_OF_WORKERS"`

	Host         string `mapstructure:"DATABASE_HOST"`
	Port         int    `mapstructure:"DATABASE_PORT"`
	User         string `mapstructure:"DATABASE_USER"`
	Password     string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName string `mapstructure:"DATABASE_NAME"`

	Users string `mapstructure:"USERS"`
}
