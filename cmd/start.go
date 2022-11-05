package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/db"
	"github.com/jbchouinard/mxremind/pkg/mail"
	"github.com/jbchouinard/mxremind/pkg/reminder"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	startCmd.Flags().Bool("migrate", false, "migrate database schema")
	viper.BindPFlag("database.migrate", rootCmd.PersistentFlags().Lookup("migrate"))

	viper.SetDefault("send_interval", 60)
	viper.SetDefault("fetch_interval", 60)

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the mail reminder service",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()

		db.SetDatabaseUrl(conf.Database.URL)
		dbpool, err := db.Connect()
		if err != nil {
			log.Fatal().Err(err).Msg("error connecting to database")
		}
		defer db.Close()

		if conf.Database.Migrate {
			if err := db.Migrate(context.Background(), conf.Database.URL); err != nil {
				log.Fatal().Err(err).Msg("error applying database migrations")
			}
		}

		// Receive and save new reminders
		fetchDone := make(chan chan<- bool)
		fetcher, messages, fetcherErrors := mail.NewMailFetcher(conf, 10, fetchDone)
		converter, reminders, converterErrors := reminder.NewReminderMailConverter(messages)
		saver, saverErrors := reminder.NewReminderSaver(dbpool, reminders)

		// Query and send due reminders
		queryDone := make(chan chan<- bool)
		querier, dueReminders, querierErrors := reminder.NewDueReminderQuerier(dbpool, queryDone)
		sender, senderErrors := reminder.NewReminderSender(dueReminders, &mail.SmtpSender{Conf: conf.SMTP})

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go fetcher.Run(time.Duration(conf.FetchInterval) * time.Second)
		go converter.Run()
		go saver.Run()
		go querier.Run(time.Duration(conf.SendInterval) * time.Second)
		go sender.Run()
		for {
			select {
			case sig := <-sigs:
				log.Info().Msgf("received signal %s, shutting down", sig)
				done := make(chan bool)
				fetchDone <- done
				queryDone <- done
				close(fetchDone)
				close(queryDone)
				<-done
				<-done
				os.Exit(0)
			case err, ok := <-fetcherErrors:
				if ok {
					log.Error().Err(err).Msg("fetcher")
				} else {
					fetcherErrors = nil
				}
			case err, ok := <-converterErrors:
				if ok {
					log.Error().Err(err).Msg("converter")
				} else {
					converterErrors = nil
				}
			case err, ok := <-saverErrors:
				if ok {
					log.Error().Err(err).Msg("saver")
				} else {
					saverErrors = nil
				}
			case err, ok := <-querierErrors:
				if ok {
					log.Error().Err(err).Msg("querier")
				} else {
					querierErrors = nil
				}
			case err, ok := <-senderErrors:
				if ok {
					log.Error().Err(err).Msg("sender")
				} else {
					senderErrors = nil
				}
			}
		}
	},
}
