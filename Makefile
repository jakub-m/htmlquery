gofiles = $(shell find . -type f -name \*.go)

default: fmt test

fmt: $(gofiles)
	go fmt ./...

test: $(gofiles)
	go test ./...

