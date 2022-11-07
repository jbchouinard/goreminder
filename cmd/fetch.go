package cmd

import (
	"fmt"
	"os"

	"github.com/emersion/go-imap"
	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/mail"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var fetchCount uint32
var trashFetchedMessages bool

func init() {
	fetchCmd.Flags().Uint32Var(&fetchCount, "count", 1, "number of messages to fetch")
	fetchCmd.Flags().BoolVar(&trashFetchedMessages, "trash", false, "trash fetched messaged")

	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch messages from the configured IMAP account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetServerConfig("imap")
		fmt.Printf("Mail for %q\n", conf.Address)
		imapClient, err := mail.ConnectImap(conf)
		if err != nil {
			log.Fatal().Err(err).Msg("")
			os.Exit(1)
		}
		defer imapClient.Logout()

		mbox, err := imapClient.Select(args[0], false)
		if err != nil {
			log.Fatal().Err(err).Msg("")
			os.Exit(1)
		}
		if mbox.Messages == 0 {
			os.Exit(0)
		}
		from, to := mail.RangeLastN(fetchCount, mbox.Messages)
		seqset := mail.RangeSeq(from, to)
		messages := make(chan *imap.Message, fetchCount)
		done := make(chan error, 1)
		go func() {
			go func() {
				if err := imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
					done <- err
				} else {
					if trashFetchedMessages {
						done <- imapClient.Move(seqset, "Trash")
					} else {
						done <- nil
					}
				}
			}()
		}()
		for message := range messages {
			fmt.Printf("%q %q\n", message.Envelope.From[0].Address(), message.Envelope.Subject)
		}
		if err := <-done; err != nil {
			log.Fatal().Err(err).Msg("")
			os.Exit(1)
		}
		os.Exit(0)
	}}
