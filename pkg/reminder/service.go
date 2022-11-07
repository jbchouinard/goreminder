package reminder

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/jbchouinard/mxremind/pkg/db"
	"github.com/jbchouinard/mxremind/pkg/mail"
)

type Component interface {
	Run()
	RunOnce()
	Close()
}

type ComponentError struct {
	Name string
	Err  error
}

func (err *ComponentError) Error() string {
	return fmt.Sprintf("%s: %s", err.Name, err.Err)
}

func (err *ComponentError) Unwrap() error {
	return err.Err
}

type Service struct {
	conf      *config.Config
	dbpool    *pgxpool.Pool
	fetcher   Component
	converter Component
	saver     Component
	querier   Component
	sender    Component
	dones     []chan<- bool
	errors    chan error
}

func NewService(ctx context.Context, conf *config.Config) (*Service, error) {
	dbpool, err := db.NewPool(ctx, conf.Database.URL)
	if err != nil {
		return nil, err
	}

	errors := make(chan error)

	// Receive and save new reminders
	fetchDone := make(chan bool)
	fetcher, messages, fetcherErrors := mail.NewMailFetcher(conf, 10, fetchDone)
	converter, reminders, converterErrors := NewReminderMailConverter(messages)
	saver, saverErrors := NewReminderSaver(dbpool, reminders)

	// Query and send due reminders
	queryDone := make(chan bool)
	querier, dueReminders, querierErrors := NewDueReminderQuerier(
		time.Duration(conf.SendInterval)*time.Second, dbpool, queryDone,
	)
	sender, senderErrors := NewReminderSender(dueReminders, &mail.SmtpSender{Conf: conf.SMTP})

	var wg sync.WaitGroup
	wg.Add(5)
	go errorPipe("fetcher", fetcherErrors, errors, &wg)
	go errorPipe("converter", converterErrors, errors, &wg)
	go errorPipe("saver", saverErrors, errors, &wg)
	go errorPipe("querier", querierErrors, errors, &wg)
	go errorPipe("sender", senderErrors, errors, &wg)
	go func(wg *sync.WaitGroup) {
		defer close(errors)
		wg.Wait()
	}(&wg)

	return &Service{
		conf:      conf,
		dbpool:    dbpool,
		saver:     saver,
		querier:   querier,
		sender:    sender,
		fetcher:   fetcher,
		converter: converter,
		dones:     []chan<- bool{fetchDone, queryDone},
		errors:    errors,
	}, nil
}

func (s *Service) Errors() <-chan error {
	return s.errors
}

func (s *Service) Start() {
	go s.querier.Run()
	go s.sender.Run()
	go s.fetcher.Run()
	go s.converter.Run()
	go s.saver.Run()
}

func (s *Service) Stop() {
	for _, done := range s.dones {
		done <- true
	}
}

func (s *Service) Drain() []error {
	errors := make([]error, 0)
	for {
		err, ok := <-s.Errors()
		if !ok {
			return errors
		}
		errors = append(errors, err)
	}
}

func (s *Service) RunOnce() {
	go func() {
		go s.converter.Run()
		go s.saver.Run()
		go s.sender.Run()
		s.fetcher.RunOnce()
		s.fetcher.Close()
		s.querier.RunOnce()
		s.querier.Close()
	}()
}

func (s *Service) Close() {
	for _, c := range s.dones {
		defer close(c)
	}
	defer s.dbpool.Close()
}

func errorPipe(name string, from <-chan error, to chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err, ok := <-from
		if !ok {
			return
		}
		to <- &ComponentError{Name: name, Err: err}
	}
}
