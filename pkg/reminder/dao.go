package reminder

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type ReminderDAO struct {
	Tx      pgx.Tx
	Context context.Context
}

// id UUID PRIMARY KEY,
// generated_from_id TEXT,
// recipient TEXT,
// content TEXT,
// due_time TIMESTAMP,
// is_sent BOOLEAN

func (dao *ReminderDAO) Scan(row pgx.Row) (*Reminder, error) {
	var rem Reminder
	err := row.Scan(
		&rem.Id,
		&rem.GeneratedById,
		&rem.Recipient,
		&rem.Content,
		&rem.DueTime,
		&rem.IsSent,
	)
	return &rem, err
}

func (dao *ReminderDAO) Load(id uuid.UUID) (*Reminder, error) {
	return dao.Scan(dao.Tx.QueryRow(
		dao.Context,
		`SELECT id, generated_from_id, recipient, content, due_time, is_sent
			FROM reminders
			WHERE id=$1`,
		id,
	))
}

func (dao *ReminderDAO) Save(rem *Reminder) error {
	_, err := dao.Tx.Exec(
		dao.Context,
		`INSERT INTO reminders
			(id, generated_from_id, recipient, content, due_time, is_sent)
			VALUES
			($1, $2, $3, $4, $5, $6)`,
		rem.Id,
		rem.GeneratedById,
		rem.Recipient,
		rem.Content,
		rem.DueTime.UTC(),
		rem.IsSent,
	)
	return err
}

func (dao *ReminderDAO) Update(rem *Reminder) error {
	_, err := dao.Tx.Exec(
		dao.Context,
		`UPDATE reminders
			SET generated_from_id = $2,
				recipient = $3,
				content = $4,
				due_time = $5,
				is_sent = $6
			WHERE id = $1`,
		rem.Id,
		rem.GeneratedById,
		rem.Recipient,
		rem.Content,
		rem.DueTime.UTC(),
		rem.IsSent,
	)
	return err
}

func (dao *ReminderDAO) Delete(rem *Reminder) error {
	_, err := dao.Tx.Exec(
		dao.Context,
		"DELETE FROM reminders WHERE id=$1",
		rem.Id,
	)
	return err
}

func (dao *ReminderDAO) QueryDue(asOf time.Time) ([]*Reminder, error) {
	now := time.Now().UTC()
	rows, err := dao.Tx.Query(
		dao.Context,
		`SELECT id, generated_from_id, recipient, content, due_time, is_sent
			FROM reminders
			WHERE due_time <= $1
			  AND NOT is_sent`,
		now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reminders := make([]*Reminder, 0)
	for rows.Next() {
		rem, err := dao.Scan(rows)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, rem)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return reminders, nil
}
