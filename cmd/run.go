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
	Use:   "run",
	Short: "Run the mail reminder",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		if migrateDatabase {
			if err := db.Migrate(context.Background(), conf.Database.URL); err != nil {
				log.Fatal().Err(err).Msg("error applying database migrations")
			}
		}
		ok := reminder.RunOnce(conf)
		if ok {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	},
}
