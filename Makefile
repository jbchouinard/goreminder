default: bin/mxremind

bin/mxremind: main.go
	go build -o $@ $<

clean:
	rm -rf bin/*

.PHONY: default clean
