GO_SRC_FILES = $(shell find ./ -type f -name '*.go')

default: bin/mxremind

bin/mxremind: $(GO_SRC_FILES)
	go build -o $@ $<

test: bin/mxremind
	docker-compose -f test/integration/docker-compose.yaml up -d
	tush-check test/integration/test_*

clean:
	rm -f bin/mxremind

.PHONY: default test clean
