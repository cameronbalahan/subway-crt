.PHONY: all build build-arm64 test run clean

BINARY=subway-crt

all: test build

build:
	go build -o $(BINARY) cmd/subway-crt/main.go

build-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/$(BINARY)-arm64 cmd/subway-crt/main.go

test:
	go test -v ./...

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -rf bin/
