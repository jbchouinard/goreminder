package cmd

import (
	"log"
	"time"

	"github.com/jbchouinard/goreminder/pkg/mail"
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
		mailConf, err := mail.MailConfigFromEnv(configEnvPrefix)
		if err != nil {
			log.Fatal(err)
		}
		cMail := make(chan *mail.Mail)
		cError := make(chan error)

		fetcher := mail.MailFetcher{Conf: mailConf, Mail: cMail, Errors: cError}
		go fetcher.Run(fetchWaitSeconds*time.Second, 10)

		for {
			select {
			case err := <-cError:
				log.Fatalf("Error: %q\n", err)
			case msg := <-cMail:
				log.Printf("%s: %s", msg.From, msg.Subject)
			}
		}
	}}
