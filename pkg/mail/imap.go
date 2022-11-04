package mail

import (
	"fmt"

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
