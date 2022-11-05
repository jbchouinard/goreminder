package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jbchouinard/mxremind/pkg/config"
)

type SmtpClient struct {
	username   string
	smtpClient *smtp.Client
}

func ConnectSmtp(conf *config.ServerConfig) (*SmtpClient, error) {
	client, err := smtp.Dial(fmt.Sprintf("%v:%v", conf.Host, conf.Port))
	if err != nil {
		return nil, err
	}

	if conf.Tls.Enabled {
		err = client.StartTLS(conf.TlsConfig())
		if err != nil {
			return nil, err
		}
	}

	if conf.Authenticated {
		err = client.Auth(smtp.PlainAuth("", conf.Address, conf.Password, conf.Host))
		if err != nil {
			return nil, err
		}
	}

	return &SmtpClient{conf.Address, client}, nil
}

func MakeMessage(from *string, to *string, subject *string, body *string) string {
	return fmt.Sprintf(
		"From: %v\r\n"+
			"To: %v\r\n"+
			"Subject: %v\r\n"+
			"\r\n"+
			"%v\r\n",
		*from,
		*to,
		*subject,
		*body,
	)
}

func (client *SmtpClient) Send(to string, subject string, body string) error {
	message := MakeMessage(&client.username, &to, &subject, &body)
	err := client.smtpClient.Mail(client.username)
	if err != nil {
		return err
	}
	err = client.smtpClient.Rcpt(to)
	if err != nil {
		return err
	}
	w, err := client.smtpClient.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (client *SmtpClient) Quit() error {
	return client.smtpClient.Quit()
}

// SmtpSender uses a new connection for every e-mail.
type SmtpSender struct {
	Conf *config.ServerConfig
}

func (s *SmtpSender) Send(to string, subject string, body string) error {
	client, err := ConnectSmtp(s.Conf)
	if err != nil {
		return err
	}
	defer client.Quit()
	return client.Send(to, subject, body)
}
