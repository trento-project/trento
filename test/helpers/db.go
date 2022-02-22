package helpers

import (
	"context"
	"os"
	"os/signal"
	"syscall"
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
		Port:     viper.GetInt("db-port"),
		User:     viper.GetString("db-user"),
		Password: viper.GetString("db-password"),
		DBName:   viper.GetString("db-name"),
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	db, err := db.InitDB(ctx, dbConfig)
	if err != nil {
		t.Fatal("could not open test database connection")
	}

	return db
}
