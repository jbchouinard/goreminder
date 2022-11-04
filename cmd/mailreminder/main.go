package main

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/jbchouinard/goreminder/pkg/mail"
)

type Mail struct {
	MessageId string
	From      string
	Subject   string
}

type MailFetcher struct {
	Conf   *mail.MailConfig
	Mail   chan<- *Mail
	Errors chan<- error
}

func seqLastN(n uint32, total uint32) *imap.SeqSet {
	seqset := new(imap.SeqSet)
	from := uint32(1)
	to := total
	if (total + 1) > n {
		from = total + 1 - n
	}
	seqset.AddRange(from, to)
	return seqset
}

func (mc *MailFetcher) Run(wait time.Duration, maxMessages uint32) error {
	imapClient, err := mail.ImapConnect(mc.Conf)
	if err != nil {
		return err
	}
	defer imapClient.Logout()
	for {
		time.Sleep(wait)
		mbox, err := imapClient.Select(mc.Conf.MailboxIn, true)
		if err != nil {
			mc.Errors <- err
			continue
		}
		seqset := seqLastN(maxMessages, mbox.Messages)
		messages := make(chan *imap.Message, maxMessages)
		done := make(chan error, 1)
		go func() {
			done <- imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()
		for message := range messages {
			if len(message.Envelope.From) < 1 {
				mc.Errors <- fmt.Errorf("message %q has no From", message.Envelope.MessageId)
				continue
			}
			mc.Mail <- &Mail{
				From:      message.Envelope.From[0].Address(),
				Subject:   message.Envelope.Subject,
				MessageId: message.Envelope.MessageId,
			}
		}
		if err := <-done; err != nil {
			mc.Errors <- err
		}

	}
}

func main() {
	mailConf, err := mail.MailConfigFromEnv("MAILREMINDER_")
	if err != nil {
		log.Fatal(err)
	}
	cMail := make(chan *Mail)
	cError := make(chan error)

	fetcher := MailFetcher{mailConf, cMail, cError}
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
