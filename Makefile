EXE_EXT ?= ""

SENDMAIL_BIN = sendmail${EXE_EXT}

default: bin/$(SENDMAIL_BIN)

bin/$(SENDMAIL_BIN):
	go build -o $@ cmd/sendmail/main.go

clean:
	rm -rf bin/*

.PHONY: default clean
