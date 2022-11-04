package mail

import (
	"crypto/tls"
	"strconv"

	"github.com/jbchouinard/goreminder/pkg/env"
)

type MailConfig struct {
	Username         string
	Password         string
	MailboxIn        string
	MailboxProcessed string
	SmtpHost         string
	SmtpPort         uint16
	SmtpTlsConfig    *tls.Config
	ImapHost         string
	ImapPort         uint16
	ImapTlsConfig    *tls.Config
}

func MailConfigFromEnv(prefix string) (*MailConfig, error) {
	env := env.EnvGetter{prefix}
	username, err := env.Get("USERNAME", nil)
	if err != nil {
		return nil, err
	}
	password, err := env.Get("PASSWORD", nil)
	if err != nil {
		return nil, err
	}
	smtpHost, err := env.Get("SMTP_HOST", nil)
	if err != nil {
		return nil, err
	}
	defaultSmtpPort := "993"
	smtpPortStr, err := env.Get("SMTP_PORT", &defaultSmtpPort)
	if err != nil {
		return nil, err
	}
	smtpPort, err := strconv.ParseUint(smtpPortStr, 10, 16)
	if err != nil {
		return nil, err
	}
	imapHost, err := env.Get("IMAP_HOST", nil)
	if err != nil {
		return nil, err
	}
	imapPortStr, err := env.Get("IMAP_PORT", nil)
	if err != nil {
		return nil, err
	}
	imapPort, err := strconv.ParseUint(imapPortStr, 10, 16)
	if err != nil {
		return nil, err
	}
	mailboxIn, err := env.Get("MAILBOX_IN", nil)
	if err != nil {
		return nil, err
	}
	mailboxProcessed, err := env.Get("MAILBOX_PROCESSED", nil)
	if err != nil {
		return nil, err
	}
	return &MailConfig{
		username,
		password,
		mailboxIn,
		mailboxProcessed,
		smtpHost,
		uint16(smtpPort),
		TlsConfig(smtpHost),
		imapHost,
		uint16(imapPort),
		TlsConfig(imapHost),
	}, nil
}

func TlsConfig(host string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
}
