all: build

.PHONY: test
test:
	go test ./internal/app/client

.PHONY: build
build:
	go build -o bin/client.bin cmd/client/main.go

.PHONY: clean
clean:
	rm -f bin/client.bin