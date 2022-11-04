package cmd

import (
	"log"

	"github.com/jbchouinard/goreminder/pkg/mail"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send <address> <subject>",
	Args:  cobra.ExactArgs(2),
	Short: "Send an email.",
	Run: func(cmd *cobra.Command, args []string) {
		mailConf, err := mail.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}
		smtpClient, err := mail.ConnectSmtp(mailConf)
		if err != nil {
			log.Fatal(err)
		}
		defer smtpClient.Quit()
		err = smtpClient.Send(args[0], args[1], "")
		if err != nil {
			log.Fatal(err)
		}
	}}
