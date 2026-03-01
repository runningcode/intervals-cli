BINARY := intervals-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test clean

build:
	go build -ldflags "-X github.com/runningcode/intervals-cli/cmd.version=$(VERSION)" -o $(BINARY) .

test:
	go test ./... -v

clean:
	rm -f $(BINARY)
