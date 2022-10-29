package mail

import (
	"strconv"

	"github.com/jbchouinard/goreminder/pkg/env"
)

type SmtpConfig struct {
	Username string
	Password string
	Host     string
	Port     uint16
}

func SmtpConfigFromEnv() (*SmtpConfig, error) {
	username, err := env.Get("SMTP_USERNAME")
	if err != nil {
		return nil, err
	}
	password, err := env.Get("SMTP_PASSWORD")
	if err != nil {
		return nil, err
	}
	host, err := env.Get("SMTP_HOST")
	if err != nil {
		return nil, err
	}
	portstr, err := env.Get("SMTP_PORT")
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseUint(portstr, 10, 16)
	if err != nil {
		return nil, err
	}
	return &SmtpConfig{username, password, host, uint16(port)}, nil
}
