#!/usr/bin/env bash
set -euo pipefail

# subway-crt Cross-Compiler for Raspberry Pi (ARM64 and ARMv7/32-bit)
#
# If you get "Exec format error" when running subway-crt-arm64, your
# Raspberry Pi is running a 32-bit OS kernel (e.g., standard 32-bit Raspberry Pi OS / armv7l).
# In that case, use bin/subway-crt-armv7.
#
# Check your Pi's kernel architecture by running: uname -m
# - If output is "aarch64", run: ./subway-crt-arm64
# - If output is "armv7l" or "armv6l", run: ./subway-crt-armv7

mkdir -p bin

echo "==> [1/2] Building 64-bit ARM binary (Raspberry Pi OS 64-bit / aarch64)..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/subway-crt-arm64 cmd/subway-crt/main.go
echo "    -> bin/subway-crt-arm64 ready"

echo "==> [2/2] Building 32-bit ARM binary (Raspberry Pi OS 32-bit / armv7l)..."
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/subway-crt-armv7 cmd/subway-crt/main.go
echo "    -> bin/subway-crt-armv7 ready"

echo ""
echo "Build complete! Output binaries:"
ls -lh bin/subway-crt-arm*
