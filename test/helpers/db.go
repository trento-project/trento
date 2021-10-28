package helpers

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDatabase(t *testing.T) *gorm.DB {
	testEnabled := viper.GetBool("db-integration-tests")
	if !testEnabled {
		t.SkipNow()
	}

	host := viper.GetString("db-host")
	port := viper.GetString("db-port")
	user := viper.GetString("db-user")
	password := viper.GetString("db-password")
	dbName := viper.GetString("db-name")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("could not open test database connection")
	}

	return db
}
