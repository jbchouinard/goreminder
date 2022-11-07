package reminder

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/db"
	"github.com/jbchouinard/mxremind/pkg/mail"
	"github.com/rs/zerolog/log"
)

func recvErrors(name string, errors <-chan error, ok *bool) {
	select {
	case err := <-errors:
		if err != nil {
			log.Error().Err(err).Msg(name)
			*ok = false
		}
	default:
	}
}

type Component interface {
	Close()
	RunOnce()
}

func runOnce(c Component, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.Close()
		c.RunOnce()
	}()
}

func RunOnce(conf *config.Config) (ok bool) {
	db.SetDatabaseUrl(conf.Database.URL)
	dbpool, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}
	defer db.Close()

	// Receive and save new reminders
	fetchDone := make(chan chan<- bool)
	fetcher, messages, fetcherErrors := mail.NewMailFetcher(conf, 10, fetchDone)
	converter, reminders, converterErrors := NewReminderMailConverter(messages)
	saver, saverErrors := NewReminderSaver(dbpool, reminders)

	// Query and send due reminders
	queryDone := make(chan chan<- bool)
	querier, dueReminders, querierErrors := NewDueReminderQuerier(dbpool, queryDone)
	sender, senderErrors := NewReminderSender(dueReminders, &mail.SmtpSender{Conf: conf.SMTP})

	var wg sync.WaitGroup
	runOnce(querier, &wg)
	runOnce(fetcher, &wg)
	runOnce(converter, &wg)
	runOnce(saver, &wg)
	runOnce(sender, &wg)
	wg.Wait()
	ok = true
	recvErrors("fetcher", fetcherErrors, &ok)
	recvErrors("querier", querierErrors, &ok)
	recvErrors("converter", converterErrors, &ok)
	recvErrors("saver", saverErrors, &ok)
	recvErrors("sender", senderErrors, &ok)
	return
}

func RunForever(conf *config.Config) {
	db.SetDatabaseUrl(conf.Database.URL)
	dbpool, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}
	defer db.Close()
	// Receive and save new reminders
	fetchDone := make(chan chan<- bool)
	fetcher, messages, fetcherErrors := mail.NewMailFetcher(conf, 10, fetchDone)
	converter, reminders, converterErrors := NewReminderMailConverter(messages)
	saver, saverErrors := NewReminderSaver(dbpool, reminders)

	// Query and send due reminders
	queryDone := make(chan chan<- bool)
	querier, dueReminders, querierErrors := NewDueReminderQuerier(dbpool, queryDone)
	sender, senderErrors := NewReminderSender(dueReminders, &mail.SmtpSender{Conf: conf.SMTP})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go fetcher.Run(time.Duration(conf.FetchInterval) * time.Second)
	go converter.Run()
	go saver.Run()
	go querier.Run(time.Duration(conf.SendInterval) * time.Second)
	go sender.Run()
	for {
		select {
		case sig := <-sigs:
			log.Info().Msgf("received signal %s, shutting down", sig)
			done := make(chan bool)
			fetchDone <- done
			queryDone <- done
			close(fetchDone)
			close(queryDone)
			<-done
			<-done
			os.Exit(0)
		case err, ok := <-fetcherErrors:
			if ok {
				log.Error().Err(err).Msg("fetcher")
			} else {
				fetcherErrors = nil
			}
		case err, ok := <-converterErrors:
			if ok {
				log.Error().Err(err).Msg("converter")
			} else {
				converterErrors = nil
			}
		case err, ok := <-saverErrors:
			if ok {
				log.Error().Err(err).Msg("saver")
			} else {
				saverErrors = nil
			}
		case err, ok := <-querierErrors:
			if ok {
				log.Error().Err(err).Msg("querier")
			} else {
				querierErrors = nil
			}
		case err, ok := <-senderErrors:
			if ok {
				log.Error().Err(err).Msg("sender")
			} else {
				senderErrors = nil
			}
		}
	}
}
