package helpers

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/trento-project/trento/internal/db"
	"gorm.io/gorm"
)

func SetupTestDatabase(t *testing.T) *gorm.DB {
	testEnabled := viper.GetBool("db-integration-tests")
	if !testEnabled {
		t.SkipNow()
	}

	dbConfig := &db.Config{
		Host:     viper.GetString("db-host"),
		Port:     viper.GetString("db-port"),
		User:     viper.GetString("db-user"),
		Password: viper.GetString("db-password"),
		DBName:   viper.GetString("db-name"),
	}

	db, err := db.InitDB(dbConfig)
	if err != nil {
		t.Fatal("could not open test database connection")
	}

	return db
}
