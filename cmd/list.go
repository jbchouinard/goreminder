package cmd

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/jbchouinard/goreminder/pkg/mail"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List mailboxes.",
	Run: func(cmd *cobra.Command, args []string) {
		mailConf, err := mail.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}
		imapClient, err := mail.ImapConnect(mailConf)
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
