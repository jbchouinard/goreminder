# MxRemind

MxRemind is a service for setting and receiving reminders via e-mail.

MxRemind queries a mailbox using IMAP to fetch e-mails that set new reminders.
Any e-mail in the mailbox that matches the reminder format will set a reminder for the sender.

## Installation

### Requirements
- golang 1.19+

```sh
go install github.com/jbchouinard/mxremind@latest
```

Or use docker-compose:

```yaml
version: '3'

services:
  db:
    image: postgres:15-alpine
    restart: always
    environment:
      - POSTGRES_USER=mxremind
      - POSTGRES_PASSWORD=mxremind
      - POSTGRES_DB=mxremind
    volumes:
      - pgdata:/var/lib/postgresql/data 
    healthcheck:
      test: [ "CMD", "pg_isready", "-q", "-d", "mxremind", "-U", "mxremind" ]
      timeout: 30s
      interval: 10s
      retries: 10

  mxremind:
    # build: https://github.com/jbchouinard/mxremind.git
    image: jbchouinard/mxremind:latest
    depends_on:
      db:
        condition: service_healthy
    environment:
      - MXREMIND_TIMEZONE=America/New York
      - MXREMIND_DATABASE_URL=postgresql://mxremind:mxremind@db/mxremind
      - MXREMIND_MAILBOX_IN=INBOX
      - MXREMIND_MAILBOX_PROCESSED=
      - MXREMIND_SMTP_ADDRESS=myname@example.com
      - MXREMIND_SMTP_PASSWORD=mypass
      - MXREMIND_SMTP_HOST=smtp.example.com
      - MXREMIND_SMTP_PORT=587
      - MXREMIND_IMAP_ADDRESS=myname@example.com
      - MXREMIND_IMAP_PASSWORD=mypass
      - MXREMIND_IMAP_HOST=imap.example.com
      - MXREMIND_IMAP_PORT=993

volumes:
  pgdata: {}
```

## Configuration

MxRemind requires a Postgresql database.

By default, MxRemind looks for the configuration file `mxremind.yaml` in the current directory.
A different configuration file can be specified with the `--config` flag.

Alternatively, options can be provided by environment variables, or CLI flags.

CLI flags have highest precedence, followed by environment variables, and finally configuration file.

| Environment Variable       | Example                             | Details                               |
|----------------------------|-------------------------------------|-------------------------------------- |
| MXREMIND_TIMEZONE          | America/Montreal                    | Default timezone for reminders.       |
| MXREMIND_FETCH_INTERVAL    | 60                                  | Interval in seconds to fetch emails.  |
| MXREMIND_SEND_INTERVAL     | 60                                  | Interval in seconds to send emails.   |
| MXREMIND_DATABASE_URL      | postgresql://user:pass@host:5432/db | Database connection string.           |
| MXREMIND_MAILBOX_IN        | INBOX                               | Mailbox to fetch from.                |
| MXREMIND_MAILBOX_PROCESSED | Trash                               | Mailbox to move processed emails.     |
| MXREMIND_SMTP_ADDRESS      | myname@example.com                  | SMTP server username.                 |
| MXREMIND_SMTP_PASSWORD     | mypassword123!                      | SMTP server password.                 |
| MXREMIND_SMTP_HOST         | smtp.example.com                    | SMTP server host.                     |
| MXREMIND_SMTP_PORT         | 587                                 | SMTP server port.                     |
| MXREMIND_IMAP_ADDRESS      | myname@example.com                  | IMAP server username.                 |
| MXREMIND_IMAP_PASSWORD     | mypassword123!                      | IMAP server password.                 |
| MXREMIND_IMAP_HOST         | imap.example.com                    | IMAP server host.                     |
| MXREMIND_IMAP_PORT         | 993                                 | IMAP server port.                     |

## Usage

Copy `mxremind.example.yaml` to `mxremind.yaml` and edit it.

Start the server:

```sh
mxremind run
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

## Tests

This repo contains integrations tests that use [tush](https://github.com/darius/tush).
Add `bin/` and `test/integration/scripts` to `PATH`, then run:

```sh
docker-compose -f test/integration/docker-compose-yaml up -d
make test
```

## License

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
