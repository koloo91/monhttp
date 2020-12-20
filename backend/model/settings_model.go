package model

type SettingsVo struct {
	DatabaseHost     string `json:"databaseHost"`
	DatabasePort     int    `json:"databasePort"`
	DatabaseUser     string `json:"databaseUser"`
	DatabasePassword string `json:"databasePassword"`
	DatabaseName     string `json:"databaseName"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}
