package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SmtpClient struct {
	username   string
	smtpClient *smtp.Client
}

func SmtpConnect(host string, port uint16, username string, password string) (*SmtpClient, error) {
	client, err := smtp.Dial(fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return nil, err
	}

	err = client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		return nil, err
	}

	err = client.Auth(smtp.PlainAuth("", username, password, host))
	if err != nil {
		return nil, err
	}

	return &SmtpClient{username, client}, nil
}

func SmtpConnectFromEnv() (*SmtpClient, error) {
	conf, err := SmtpConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return SmtpConnect(conf.Host, conf.Port, conf.Username, conf.Password)
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
