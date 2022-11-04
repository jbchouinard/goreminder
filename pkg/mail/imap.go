package mail

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func ImapConnect(conf *MailConfig) (*client.Client, error) {
	client, err := client.DialTLS(fmt.Sprintf("%v:%v", conf.ImapHost, conf.ImapPort), conf.ImapTlsConfig)
	if err != nil {
		return nil, err
	}
	if err := client.Login(conf.Username, conf.Password); err != nil {
		return nil, err
	}
	return client, nil
}

type Mail struct {
	MessageId string
	From      string
	Subject   string
}

type MailFetcher struct {
	Conf   *MailConfig
	Mail   chan<- *Mail
	Errors chan<- error
}

func (mc *MailFetcher) Run(wait time.Duration, maxMessages uint32) error {
	imapClient, err := ImapConnect(mc.Conf)
	if err != nil {
		return err
	}
	defer imapClient.Logout()
	for {
		time.Sleep(wait)
		mbox, err := imapClient.Select(mc.Conf.MailboxIn, false)
		if err != nil {
			mc.Errors <- err
			continue
		}
		mc.Conf.Logf("Contains %d messages", mbox.Messages)
		if mbox.Messages == 0 {
			continue
		}
		from, to := rangeLastN(maxMessages, mbox.Messages)
		mc.Conf.Logf("Fetching messages %d..%d", from, to)
		seqset := rangeSeq(from, to)
		messages := make(chan *imap.Message, maxMessages)
		done := make(chan error, 1)
		go func() {
			if err := imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
				done <- err
			} else {
				done <- imapClient.Move(seqset, mc.Conf.MailboxProcessed)
			}
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

type MailFetchError struct {
	Conf MailConfig
	Err  error
}

func (mfe *MailFetchError) Error() string {
	return fmt.Sprintf("%s: %s", mfe.Conf.Describe(), mfe.Err)
}

func (mfe *MailFetchError) Unwrap() error {
	return mfe.Err
}

func rangeLastN(n uint32, total uint32) (uint32, uint32) {
	from := uint32(1)
	to := total
	if (total + 1) > n {
		from = total + 1 - n
	}
	return from, to
}

func rangeSeq(from uint32, to uint32) *imap.SeqSet {
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	return seqset
}
