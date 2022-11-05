package cmd

import (
	"context"
	"log"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetDatabaseConfig("database")
		if err := db.Migrate(context.Background(), conf.URL); err != nil {
			log.Fatal(err)
		}
	}}
