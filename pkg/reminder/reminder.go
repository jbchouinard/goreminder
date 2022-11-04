package reminder

import (
	"time"
)

type Sender interface {
	Send(to string, subject string, body string) error
}

type Reminder struct {
	AtTime    time.Time
	ToAddress string
	Message   string
}

func (rem *Reminder) SendWith(sender Sender) error {
	return sender.Send(rem.ToAddress, rem.Message, "")
}
