package main

import (
	"log"

	"github.com/jbchouinard/goreminder/pkg/mail"
)

func main() {
	mailConf, err := mail.MailConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	smtpClient, err := mail.SmtpConnect(mailConf)
	if err != nil {
		log.Fatal(err)
	}
	defer smtpClient.Quit()
	err = smtpClient.Send("me@jbchouinard.net", "Hello there!", "Did you get this?")
	if err != nil {
		log.Fatal(err)
	}
}
