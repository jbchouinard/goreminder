package reminder

import (
	"time"
)

type Sender interface {
	Send(to string, subject string, body string) error
}

type Reminder struct {
	AtTime  time.Time
	To      string
	Message string
}

func (rem *Reminder) SendWith(sender Sender) error {
	return sender.Send(rem.To, rem.Message, "")
}

func Parse(s string) (*Reminder, error) {
	return nil, nil
}
