package db

import (
	"github.com/spf13/cobra"
)

func AddDBFlags(cmd *cobra.Command) {
	var dbHost string
	var dbPort int
	var dbUser string
	var dbPassword string
	var dbName string

	cmd.PersistentFlags().StringVar(&dbHost, "db-host", "localhost", "The database host")
	cmd.PersistentFlags().IntVar(&dbPort, "db-port", 5432, "The database port to connect to")
	cmd.PersistentFlags().StringVar(&dbUser, "db-user", "postgres", "The database user")
	cmd.PersistentFlags().StringVar(&dbPassword, "db-password", "postgres", "The database password")
	cmd.PersistentFlags().StringVar(&dbName, "db-name", "trento", "The database name that the application will use")
}
