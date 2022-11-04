package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbchouinard/goreminder/pkg/mail"
	"github.com/jbchouinard/goreminder/pkg/reminder"
	"github.com/spf13/cobra"
)

const fetchWaitSeconds = 10
const configEnvPrefix = "MAILREMINDER_"

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the mail reminder service.",
	Run: func(cmd *cobra.Command, args []string) {
		mailConf, err := mail.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}

		messages := make(chan *mail.Mail, 100)
		fetchErrors := make(chan error)
		fetchDone := make(chan chan<- bool)
		reminders := make(chan *reminder.Reminder, 100)
		convertErrors := make(chan error)

		fetcher := mail.MailFetcher{
			Conf:   mailConf,
			Mail:   messages,
			Errors: fetchErrors,
			Done:   fetchDone,
		}

		converter := reminder.ReminderMailConverter{
			Mail:      messages,
			Reminders: reminders,
			Errors:    convertErrors,
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go converter.Run()
		go fetcher.Run(fetchWaitSeconds*time.Second, 10)
		for {
			select {
			case sig := <-sigs:
				log.Printf("Received signal %s, shutting down", sig)
				done := make(chan bool)
				fetchDone <- done
				<-done
				os.Exit(0)
			case err := <-fetchErrors:
				log.Printf("Error: %q\n", err)
			case err := <-convertErrors:
				log.Printf("Error: %q\n", err)
			case rem := <-reminders:
				log.Printf("TO %s AT %s: %s", rem.Recipient, rem.DueTime, rem.Content)
			}
		}
	}}
