package mail

import (
	"crypto/tls"
	"strconv"

	"github.com/jbchouinard/goreminder/pkg/env"
)

type MailConfig struct {
	Username      string
	Password      string
	SmtpHost      string
	SmtpPort      uint16
	SmtpTlsConfig *tls.Config
	ImapHost      string
	ImapPort      uint16
	ImapTlsConfig *tls.Config
}

func MailConfigFromEnv() (*MailConfig, error) {
	username, err := env.Get("MAIL_USERNAME")
	if err != nil {
		return nil, err
	}
	password, err := env.Get("MAIL_PASSWORD")
	if err != nil {
		return nil, err
	}
	smtpHost, err := env.Get("MAIL_SMTP_HOST")
	if err != nil {
		return nil, err
	}
	smtpPortStr, err := env.Get("MAIL_SMTP_PORT")
	if err != nil {
		return nil, err
	}
	smtpPort, err := strconv.ParseUint(smtpPortStr, 10, 16)
	if err != nil {
		return nil, err
	}
	imapHost, err := env.Get("MAIL_IMAP_HOST")
	if err != nil {
		return nil, err
	}
	imapPortStr, err := env.Get("MAIL_IMAP_PORT")
	if err != nil {
		return nil, err
	}
	imapPort, err := strconv.ParseUint(imapPortStr, 10, 16)
	if err != nil {
		return nil, err
	}
	return &MailConfig{username, password, smtpHost, uint16(smtpPort), TlsConfig(smtpHost), imapHost, uint16(imapPort), TlsConfig(imapHost)}, nil
}

func TlsConfig(host string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
}
