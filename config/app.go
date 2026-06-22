package config

import "online-ordering-system/models"

type Config struct {
	Port      string
	DBName    string
	JWTSecret string
}

var AppConfig = Config{
	Port:      ":8080",
	DBName:    "ordering.db",
	JWTSecret: "online-ordering-secret-key-2024",
}

func InitDB() {
	models.ConnectDB(AppConfig.DBName)
}
