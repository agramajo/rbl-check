.PHONY: clean 

rbl-check: $(wildcard *.go)
	go build -ldflags="-s -w" -o $@ $<

clean:
	go clean

