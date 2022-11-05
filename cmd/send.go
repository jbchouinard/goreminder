package cmd

import (
	"log"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/mail"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send <address> <subject>",
	Args:  cobra.ExactArgs(2),
	Short: "Send an email with the configured SMTP account",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetServerConfig("smtp")
		smtpClient, err := mail.ConnectSmtp(conf)
		if err != nil {
			log.Fatal(err)
		}
		defer smtpClient.Quit()
		err = smtpClient.Send(args[0], args[1], "")
		if err != nil {
			log.Fatal(err)
		}
	}}
