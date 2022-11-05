# MxRemind

MxRemind is a service for setting and receiving reminders via e-mail.

MxRemind is not a mail server; it queries a mailbox using IMAP to fetch e-mails that set new reminders.
Any e-mail in the mailbox that matches the reminder format will set a reminder for the sender.

## Installation

Requirements
- golang 1.19+

```sh
go install github.com/jbchouinard/mxremind@latest
```

## Configuration

MxRemind requires a Postgresql database.

By default, MxRemind will read the configuration file `mxremind.yaml` in the current directory if it exists.
A different configuration file can be specified with the `--config` flag.

Alternatively, options can be provided by environment variables, or some by CLI flags.

CLI flags have highest precedence, followed by environment variables, and finally the configuration file.

| YAML Key          | Environment Variable       | Example                             | Details                                             |
|-------------------|----------------------------|-------------------------------------|-----------------------------------------------------|
| timezone          | MXREMIND_TIMEZONE          | America/Montreal                    | Default timezone for reminders.                     |
| fetch_interval    | MXREMIND_FETCH_INTERVAL    |                                     | Interval in seconds to fetch emails.                |
| send_interval     | MXREMIND_SEND_INTERVAL     |                                     | Interval in seconds to send reminder emails.        |
| database.url      | MXREMIND_DATABASE_URL      | postgresql://user:pass@host:5432/db | Database connection string.                         |
| database.migrate  | MXREMIND_DATABASE_migrate  | true                                | Automatically apply database migrations on startup. |
| mailbox.in        | MXREMIND_MAILBOX_IN        | INBOX/Reminders                     | Mailbox to monitor for setting reminders.           |
| mailbox.processed | MXREMIND_MAILBOX_PROCESSED | INBOX/Reminders/Processed           | Mailbox to move processed reminders to.             |
| smtp.address      | MXREMIND_SMTP_ADDRESS      | myname@example.com                  | SMTP server username.                               |
| smtp.password     | MXREMIND_SMTP_PASSWORD     | mypassword123!                      | SMTP server password.                               |
| smtp.host         | MXREMIND_SMTP_HOST         | smtp.example.com                    | SMTP server host.                                   |
| smtp.port         | MXREMIND_SMTP_PORT         | 587                                 | SMTP server port.                                   |
| imap.address      | MXREMIND_IMAP_ADDRESS      | myname@example.com                  | IMAP server username.                               |
| imap.password     | MXREMIND_IMAP_PASSWORD     | mypassword123!                      | IMAP server password.                               |
| imap.host         | MXREMIND_IMAP_HOST         | imap.example.com                    | IMAP server host.                                   |
| imap.port         | MXREMIND_IMAP_PORT         | 993                                 | IMAP server port.                                   |                         |

## Usage

Start the server:

```sh
mxremind start
```

A new reminder is set by sending an e-mail to the configured mailbox with a subject matching one of
these formats:

| Format                   | Example                       |
|--------------------------|-------------------------------|
| HH:MM message            | 15:04 do the thing            |
| tomorrow HH:MM message   | tomorrow 15:04 do the thing   |
| MM/DD HH:MM message      | 12/04 08:00 do the thing      |
| YYYY-MM-DD HH:MM message | 2023-04-05 12:00 do the thing |

A reminder e-mail will be sent back to the sender with the message as subject at the specified time.

# License

Copyright 2022 Jerome Boisvert-Chouinard

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
