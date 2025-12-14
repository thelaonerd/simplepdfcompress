#!/bin/bash
# Build script for SimplePDFCompress with optimization flags
# -s: Omit the symbol table and debug information
# -w: Omit the DWARF symbol table

set -e # Exit immediately if a command exits with a non-zero status

# Create bin directory
mkdir -p bin

echo "Building Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/simplepdfcompress-linux-amd64 .
cp bin/simplepdfcompress-linux-amd64 simplepdfcompress

echo "Building Windows (amd64)..."
if command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o bin/simplepdfcompress-windows-amd64.exe .
else
    echo "Skipping Windows build: MinGW compiler (x86_64-w64-mingw32-gcc) not found."
    echo "To build for Windows, install mingw-w64."
fi

echo "Building macOS..."
echo "Skipping macOS build: Cross-compiling CGO (Fyne) for macOS requires a specialized toolchain (osxcross) not present."
# GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/simplepdfcompress-darwin-arm64 .


echo "Builds completed successfully in ./bin/"
ls -lh bin/
