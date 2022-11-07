package cmd

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/mail"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List mailboxes for the configured IMAP account",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetServerConfig("imap")
		fmt.Printf("Mailboxes for %q\n", conf.Address)
		imapClient, err := mail.ConnectImap(conf)
		if err != nil {
			log.Fatal(err)
		}
		defer imapClient.Logout()

		mailboxes := make(chan *imap.MailboxInfo, 10)
		done := make(chan error, 1)
		go func() {
			done <- imapClient.List("", "*", mailboxes)
		}()

		for m := range mailboxes {
			fmt.Println(m.Name)
		}
	}}
