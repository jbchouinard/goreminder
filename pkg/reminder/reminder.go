package reminder

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jbchouinard/goreminder/pkg/mail"
)

type Sender interface {
	Send(to string, subject string, body string) error
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
		Id:            uuid.New(),
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
}

func (rmc *ReminderMailConverter) Run() {
	for {
		msg, ok := <-rmc.Mail
		if !ok {
			close(rmc.Reminders)
			close(rmc.Errors)
			return
		}
		rem, err := ReminderFromMail(msg)
		if err != nil {
			rmc.Errors <- fmt.Errorf("%s: %w", msg.MessageId, err)
		} else {
			rmc.Reminders <- rem
		}
	}
}
