package main

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/jbchouinard/goreminder/pkg/mail"
)

func main() {
	mailConf, err := mail.MailConfigFromEnv("MAIL_")
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

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}
}
