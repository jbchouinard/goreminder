package cmd

import (
	"context"
	"os"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/db"
	"github.com/jbchouinard/mxremind/pkg/reminder"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	runCmd.Flags().BoolVar(&migrateDatabase, "migrate", false, "migrate database schema")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "batch",
	Short: "Process a single batch of mail reminders",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		conf := config.GetConfig()

		if migrateDatabase {
			if err := db.Migrate(ctx, conf.Database.URL); err != nil {
				log.Fatal().Err(err).Msg("error applying database migrations")
			}
		}

		service, err := reminder.NewService(ctx, conf)
		if err != nil {
			log.Fatal().Err(err).Msg("error initializing service")
		}
		defer service.Close()

		service.RunOnce()

		errors := service.Drain()
		for _, err := range errors {
			log.Error().Err(err).Msg("")
		}

		if len(errors) > 0 {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	},
}
