#!/usr/bin/env bash
set -euo pipefail

echo "==> Cross-compiling subway-crt for Raspberry Pi 3B (ARM64)..."
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=0

mkdir -p bin
go build -ldflags="-s -w" -o bin/subway-crt-arm64 cmd/subway-crt/main.go
echo "==> Build complete: bin/subway-crt-arm64"
