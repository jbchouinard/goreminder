package cmd

import (
	"context"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/db"
	"github.com/jbchouinard/mxremind/pkg/reminder"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	startCmd.Flags().BoolVar(&migrateDatabase, "migrate", false, "migrate database schema")

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the mail reminder service",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		if migrateDatabase {
			if err := db.Migrate(context.Background(), conf.Database.URL); err != nil {
				log.Fatal().Err(err).Msg("error applying database migrations")
			}
		}
		reminder.RunForever(conf)
	},
}
