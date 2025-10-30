# Sumb Task Management Application Makefile

# Variables
BINARY_NAME=sumb
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-windows-amd64.exe .

# Install to system
.PHONY: install
install: build
	@echo "Installing ${BINARY_NAME}..."
	cp ${BINARY_NAME} /usr/local/bin/
	@echo "Installation complete! You can now use 'sumb' command."

# Install to user's home directory
.PHONY: install-user
install-user: build
	@echo "Installing ${BINARY_NAME} to user directory..."
	mkdir -p ~/.local/bin
	cp ${BINARY_NAME} ~/.local/bin/
	@echo "Installation complete! Add ~/.local/bin to your PATH if not already there."
