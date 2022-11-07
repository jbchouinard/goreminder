default: bin/mxremind

bin/mxremind: main.go
	go build -o $@ $<

test: bin/mxremind
	docker-compose -f test/integration/docker-compose.yaml up -d
	tush-check test/integration/test_*

clean:
	rm -f bin/mxremind

.PHONY: default test clean
