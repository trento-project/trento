package db

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func InitDB(ctx context.Context, config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)

	var db *gorm.DB
	var err error

	err = retry.Do(
		func() error {
			log.Info("Connecting to the database")
			// TODO: since we are dealing with eventual consistency, we can't enforce foreign key constraints in our projected models.
			// This disables foreign key constraints enforcement at global level.
			// In a future we will enable this on a per-model basis via dedicated migrations and disabling the automigration feature.
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
			})
			if err != nil {
				return err
			}

			return nil
		},
		retry.OnRetry(func(_ uint, err error) {
			log.Error(err)
		}),
		retry.Delay(1*time.Second),
		retry.MaxJitter(2*time.Second),
		retry.Attempts(8),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
	)

	return db, err
}
