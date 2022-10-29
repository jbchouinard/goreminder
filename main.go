package main

import (
	"log"

	"github.com/jbchouinard/goreminder/mail"
)

func main() {
	smtpClient, err := mail.SmtpConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	err = smtpClient.Send("me@jbchouinard.net", "Hello there!", "Did you get this?")
	if err != nil {
		log.Fatal(err)
	}
	err = smtpClient.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
