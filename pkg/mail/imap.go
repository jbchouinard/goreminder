package mail

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/rs/zerolog/log"
)

func ConnectImap(conf *config.ServerConfig) (*client.Client, error) {
	var c *client.Client
	var err error
	addr := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	if conf.Tls.Enabled {
		c, err = client.DialTLS(addr, conf.TlsConfig())
	} else {
		c, err = client.Dial(addr)
	}
	if err != nil {
		return nil, err
	}
	if conf.Authenticated {
		if err := c.Login(conf.Address, conf.Password); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Mail struct {
	MessageId string
	From      string
	Subject   string
	Location  *time.Location
}

type MailFetcher struct {
	Conf        *config.Config
	MaxMessages uint32
	Done        <-chan bool
	Mail        chan<- *Mail
	Errors      chan<- error
}

func NewMailFetcher(
	conf *config.Config, maxMessages uint32, done <-chan bool,
) (*MailFetcher, <-chan *Mail, <-chan error) {
	mail := make(chan *Mail, maxMessages)
	errors := make(chan error, 1)
	return &MailFetcher{conf, maxMessages, done, mail, errors}, mail, errors
}

func (f *MailFetcher) RunOnce() {
	imapClient, err := ConnectImap(f.Conf.IMAP)
	if err != nil {
		f.Errors <- err
		return
	}
	defer imapClient.Logout()
	mbox, err := imapClient.Select(f.Conf.Mailbox.In, false)
	if err != nil {
		f.Errors <- err
		return
	}
	log.Info().Msgf("%s/%s contains %d messages", f.Conf.IMAP.Address, f.Conf.Mailbox.In, mbox.Messages)
	if mbox.Messages == 0 {
		return
	}
	from, to := RangeLastN(f.MaxMessages, mbox.Messages)
	log.Info().Msgf("%s/%s fetching messages %d-%d", f.Conf.IMAP.Address, f.Conf.Mailbox.In, from, to)
	seqset := RangeSeq(from, to)
	messages := make(chan *imap.Message, f.MaxMessages)
	done := make(chan error, 1)
	go func() {
		if err := imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
			done <- err
		} else {
			done <- imapClient.Move(seqset, f.Conf.Mailbox.Processed)
		}
	}()
	for message := range messages {
		if len(message.Envelope.From) == 0 {
			f.Errors <- fmt.Errorf("message %q has no From", message.Envelope.MessageId)
			return
		}
		f.Mail <- &Mail{
			From:      message.Envelope.From[0].Address(),
			Subject:   message.Envelope.Subject,
			MessageId: message.Envelope.MessageId,
			Location:  f.Conf.Location(),
		}
	}
	if err := <-done; err != nil {
		f.Errors <- err
	}
}

func (f *MailFetcher) Close() {
	close(f.Mail)
	close(f.Errors)
}

func (f *MailFetcher) Run() {
	defer f.Close()
	for {
		select {
		case <-f.Done:
			return
		case <-time.After(time.Duration(f.Conf.FetchInterval) * time.Second):
		}
		f.RunOnce()
	}
}

type MailFetchError struct {
	Conf config.Config
	Err  error
}

func (mfe *MailFetchError) Error() string {
	return fmt.Sprintf("%s/%s %s", mfe.Conf.IMAP.Address, mfe.Conf.Mailbox.In, mfe.Err)
}

func (mfe *MailFetchError) Unwrap() error {
	return mfe.Err
}

func RangeLastN(n uint32, total uint32) (uint32, uint32) {
	from := uint32(1)
	to := total
	if (total + 1) > n {
		from = total + 1 - n
	}
	return from, to
}

func RangeSeq(from uint32, to uint32) *imap.SeqSet {
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	return seqset
}
