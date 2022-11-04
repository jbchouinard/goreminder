package main

import (
	"log"
	"time"

	"github.com/jbchouinard/goreminder/pkg/mail"
)

func main() {
	mailConf, err := mail.MailConfigFromEnv("MAILREMINDER_")
	if err != nil {
		log.Fatal(err)
	}
	cMail := make(chan *mail.Mail)
	cError := make(chan error)

	fetcher := mail.MailFetcher{Conf: mailConf, Mail: cMail, Errors: cError}
	go fetcher.Run(10*time.Second, 10)

	for {
		select {
		case err := <-cError:
			log.Fatalf("Error: %q\n", err)
		case msg := <-cMail:
			log.Printf("%s: %s", msg.From, msg.Subject)
		}
	}
}
