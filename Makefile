.PHONY: all build build-arm64 build-armv7 build-pi test run clean

BINARY=subway-crt

all: test build

build:
	go build -o $(BINARY) cmd/subway-crt/main.go

# 64-bit ARM (Raspberry Pi OS 64-bit / aarch64)
build-arm64:
	mkdir -p bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/$(BINARY)-arm64 cmd/subway-crt/main.go

# 32-bit ARM (Raspberry Pi OS 32-bit / armv7l) - solves "Exec format error" on 32-bit OS
build-armv7:
	mkdir -p bin
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/$(BINARY)-armv7 cmd/subway-crt/main.go

# Build both Pi architectures
build-pi: build-arm64 build-armv7

test:
	go test -v ./...

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -rf bin/
