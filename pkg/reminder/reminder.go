package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbchouinard/mxremind/pkg/mail"
	"github.com/rs/zerolog/log"
)

type Sender interface {
	Send(to string, subject string, body string) error
}

type PrintSender struct{}

func (PrintSender) Send(to string, subject string, body string) error {
	fmt.Printf("TO: %s\nSUBJECT: %s\nBODY:\n%s\n", to, subject, body)
	return nil
}

type Reminder struct {
	Id            uuid.UUID
	GeneratedById string
	DueTime       time.Time
	Recipient     string
	Content       string
	IsSent        bool
}

func (rem *Reminder) SendWith(sender Sender) error {
	return sender.Send(rem.Recipient, rem.Content, "")
}

func ReminderFromMail(m *mail.Mail) (*Reminder, error) {
	dueTime, content, err := parseSpec(m.Subject, m.Location)
	if err != nil {
		return nil, err
	}
	return &Reminder{
		Id:            uuid.Must(uuid.NewV1()),
		GeneratedById: m.MessageId,
		DueTime:       dueTime,
		Recipient:     m.From,
		Content:       content,
		IsSent:        false,
	}, nil
}

type ReminderMailConverter struct {
	Mail      <-chan *mail.Mail
	Reminders chan<- *Reminder
	Errors    chan<- error
	finished  bool
}

func (rmc *ReminderMailConverter) RunOnce() {
	msg, ok := <-rmc.Mail
	if !ok {
		rmc.finished = true
		return
	}
	rem, err := ReminderFromMail(msg)
	if err != nil {
		rmc.Errors <- fmt.Errorf("%s - %s: %w", msg.From, msg.MessageId, err)
	} else {
		rmc.Reminders <- rem
	}
	return
}

func (rmc *ReminderMailConverter) Close() {
	close(rmc.Reminders)
	close(rmc.Errors)
}

func (rmc *ReminderMailConverter) Run() {
	defer rmc.Close()
	for !rmc.finished {
		rmc.RunOnce()
	}
}

func NewReminderMailConverter(mail <-chan *mail.Mail) (*ReminderMailConverter, <-chan *Reminder, <-chan error) {
	reminders := make(chan *Reminder)
	errors := make(chan error, 1)
	return &ReminderMailConverter{mail, reminders, errors, false}, reminders, errors
}

type ReminderSaver struct {
	Pool      *pgxpool.Pool
	Reminders <-chan *Reminder
	Errors    chan<- error
	finished  bool
}

func NewReminderSaver(pool *pgxpool.Pool, reminders <-chan *Reminder) (*ReminderSaver, <-chan error) {
	errors := make(chan error, 1)
	return &ReminderSaver{pool, reminders, errors, false}, errors
}

func (rs *ReminderSaver) RunOnce() {
	rem, ok := <-rs.Reminders
	if !ok {
		rs.finished = true
		return
	}
	ctx := context.Background()
	tx, err := rs.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		rs.Errors <- err
		return
	}
	defer tx.Rollback(ctx)
	dao := ReminderDAO{Tx: tx, Context: ctx}
	if err = dao.Save(rem); err != nil {
		rs.Errors <- err
	} else {
		if err := tx.Commit(ctx); err != nil {
			rs.Errors <- err
		}
	}
	return
}

func (rs *ReminderSaver) Close() {
	close(rs.Errors)
}

func (rs *ReminderSaver) Run() {
	defer rs.Close()
	for {
		for !rs.finished {
			rs.RunOnce()
		}
	}
}

type ReminderSender struct {
	Sender    Sender
	Errors    chan<- error
	Reminders <-chan *Reminder
	finished  bool
}

func NewReminderSender(reminders <-chan *Reminder, sender Sender) (*ReminderSender, <-chan error) {
	errors := make(chan error, 1)
	return &ReminderSender{sender, errors, reminders, false}, errors
}

func (rs *ReminderSender) RunOnce() {
	rem, ok := <-rs.Reminders
	if !ok {
		rs.finished = true
		return
	}
	if err := rs.Sender.Send(
		rem.Recipient,
		"Reminder: "+rem.Content,
		"",
	); err != nil {
		rs.Errors <- err
	}
	return
}

func (rs *ReminderSender) Close() {
	close(rs.Errors)
}

func (rs *ReminderSender) Run() {
	defer rs.Close()
	for !rs.finished {
		rs.RunOnce()
	}
}

type DueReminderQuerier struct {
	Pool      *pgxpool.Pool
	Done      <-chan chan<- bool
	Reminders chan<- *Reminder
	Errors    chan<- error
}

func NewDueReminderQuerier(pool *pgxpool.Pool, done <-chan chan<- bool) (*DueReminderQuerier, <-chan *Reminder, <-chan error) {
	reminders := make(chan *Reminder)
	errors := make(chan error, 1)
	return &DueReminderQuerier{pool, done, reminders, errors}, reminders, errors
}

func (q *DueReminderQuerier) RunOnce() {
	ctx := context.Background()
	tx, err := q.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		q.Errors <- err
		return
	}
	defer tx.Rollback(ctx)
	dao := ReminderDAO{Tx: tx, Context: ctx}
	rems, err := dao.QueryDue(time.Now().UTC())
	log.Info().Msgf("found %d reminders due", len(rems))
	if err != nil {
		q.Errors <- err
		return
	}
	for _, rem := range rems {
		rem.IsSent = true
		if err := dao.Update(rem); err != nil {
			q.Errors <- err
		} else {
			if err := tx.Commit(ctx); err != nil {
				q.Errors <- err
			} else {
				q.Reminders <- rem

			}
		}
	}
}

func (q *DueReminderQuerier) Close() {
	close(q.Reminders)
	close(q.Errors)
}

func (q *DueReminderQuerier) Run(wait time.Duration) {
	defer q.Close()
	for {
		select {
		case done := <-q.Done:
			done <- true
			return
		case <-time.After(wait):
		}
		q.RunOnce()
	}
}
