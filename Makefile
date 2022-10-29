EXE_EXT ?= ""

SENDMAIL_BIN = sendmail${EXE_EXT}
LISTMAILBOXES_BIN = listmailboxes${EXE_EXT}

default: bin/$(SENDMAIL_BIN) bin/$(LISTMAILBOXES_BIN)

bin/$(SENDMAIL_BIN): cmd/sendmail/main.go
bin/$(LISTMAILBOXES_BIN): cmd/listmailboxes/main.go
bin/$(SENDMAIL_BIN) bin/$(LISTMAILBOXES_BIN):
	go build -o $@ $<

clean:
	rm -rf bin/*

.PHONY: default clean
