.PHONY: all build build-all clean install uninstall

BINARY_NAME=convrt
VERSION=1.0.0

# Default target
all: build

# Build for current platform
build:
	go build -o ${BINARY_NAME}$(shell go env GOEXE)

# Build for all platforms
build-all:
	# Windows
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}_${VERSION}_windows_amd64.exe
	GOOS=windows GOARCH=386 go build -o ${BINARY_NAME}_${VERSION}_windows_386.exe
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=386 go build -o ${BINARY_NAME}_${VERSION}_linux_386
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}_${VERSION}_linux_arm64
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}_${VERSION}_darwin_arm64

# Clean build artifacts
clean:
	rm -f ${BINARY_NAME}*

# Install based on platform
install:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File install.ps1
else
	./install.sh
endif

# Uninstall based on platform
uninstall:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File uninstall.ps1
else
	./uninstall.sh
endif
