package mail

import (
	"fmt"
	"net/smtp"
)

type SmtpClient struct {
	username   string
	smtpClient *smtp.Client
}

func ConnectSmtp(conf *MailConfig) (*SmtpClient, error) {
	client, err := smtp.Dial(fmt.Sprintf("%v:%v", conf.SmtpHost, conf.SmtpPort))
	if err != nil {
		return nil, err
	}

	err = client.StartTLS(conf.SmtpTlsConfig)
	if err != nil {
		return nil, err
	}

	err = client.Auth(smtp.PlainAuth("", conf.SmtpUsername, conf.SmtpPassword, conf.SmtpHost))
	if err != nil {
		return nil, err
	}

	return &SmtpClient{conf.SmtpUsername, client}, nil
}

func MakeMessage(from *string, to *string, subject *string, body *string) string {
	return fmt.Sprintf("From: %v\r\nTo: %v\r\nSubject: %v\r\n\r\n%v", *from, *to, *subject, *body)
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
