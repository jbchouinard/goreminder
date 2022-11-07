package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	Use:   "run",
	Short: "Run the mail reminder service",
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

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		service.Start()
		for {
			select {
			case sig := <-sigs:
				log.Info().Msgf("received signal %s, shutting down", sig)
				service.Stop()
				os.Exit(0)
			case err := <-service.Errors():
				log.Error().Err(err).Msg("")
			}
		}
	},
}
