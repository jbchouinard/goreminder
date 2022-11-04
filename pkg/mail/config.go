package mail

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type MailConfig struct {
	MailboxIn        string
	MailboxProcessed string
	SmtpUsername     string
	SmtpPassword     string
	SmtpHost         string
	SmtpPort         uint16
	SmtpTlsConfig    *tls.Config
	ImapUsername     string
	ImapPassword     string
	ImapHost         string
	ImapPort         uint16
	ImapTlsConfig    *tls.Config
}

func ReadConfig() (*MailConfig, error) {
	for _, key := range []string{
		"mailbox.in", "mailbox.processed",
		"smtp.username", "smtp.password", "smtp.host", "smtp.port",
		"imap.username", "imap.password", "imap.host", "imap.port",
	} {
		if !viper.IsSet(key) {
			return nil, fmt.Errorf("missing configuration key %q", key)
		}
	}

	smtpHost := viper.GetString("smtp.host")
	imapHost := viper.GetString("imap.host")
	return &MailConfig{
		MailboxIn:        viper.GetString("mailbox.in"),
		MailboxProcessed: viper.GetString("mailbox.processed"),
		SmtpUsername:     viper.GetString("smtp.username"),
		SmtpPassword:     viper.GetString("smtp.password"),
		SmtpHost:         smtpHost,
		SmtpPort:         viper.GetUint16("smtp.port"),
		SmtpTlsConfig:    TlsConfig(smtpHost),
		ImapUsername:     viper.GetString("imap.username"),
		ImapPassword:     viper.GetString("imap.password"),
		ImapHost:         imapHost,
		ImapPort:         viper.GetUint16("imap.port"),
		ImapTlsConfig:    TlsConfig(imapHost),
	}, nil
}

func (mc *MailConfig) Describe() string {
	return fmt.Sprintf("%s:%s", mc.ImapUsername, mc.MailboxIn)
}

func (mc *MailConfig) Log(message string) {
	log.Printf("%s: %s", mc.Describe(), message)
}

func (mc *MailConfig) Logf(format string, p ...any) {
	mc.Log(fmt.Sprintf(format, p...))
}

func TlsConfig(host string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
}
