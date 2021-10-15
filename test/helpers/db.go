package helpers

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDatabase() *gorm.DB {
	// TODO: refactor this in a common infrastructure init package

	viper.SetDefault("db-host", "localhost")
	viper.SetDefault("db-port", "32432")
	viper.SetDefault("db-user", "postgres")
	viper.SetDefault("db-password", "postgres")
	viper.SetDefault("db-name", "trento_test")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	host := viper.GetString("db-host")
	port := viper.GetString("db-port")
	user := viper.GetString("db-user")
	password := viper.GetString("db-password")
	dbName := viper.GetString("db-name")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
