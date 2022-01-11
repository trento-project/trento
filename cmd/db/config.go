package db

import (
	"github.com/spf13/viper"
	"github.com/trento-project/trento/internal/db"
)

func LoadConfig() *db.Config {
	return &db.Config{
		Host:     viper.GetString("db-host"),
		Port:     viper.GetInt("db-port"),
		User:     viper.GetString("db-user"),
		Password: viper.GetString("db-password"),
		DBName:   viper.GetString("db-name"),
	}
}
