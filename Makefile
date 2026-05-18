ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY_EXT  := .exe
    RM          := cmd /C del /Q /F
    RMDIR       := cmd /C rmdir /S /Q
    MKDIR       := cmd /C mkdir
    SEP         := \\
else
    DETECTED_OS := $(shell uname -s)
    BINARY_EXT  :=
    RM          := rm -f
    RMDIR       := rm -rf
    MKDIR       := mkdir -p
    SEP         := /
endif

BINARY := nostr-reader-api$(BINARY_EXT)
BIN_DIR := bin
GOFLAGS = -ldflags="-s -w"

.PHONY: all build run clean

all: clean build

build:
	@echo Building for $(DETECTED_OS)...
	go build $(GOFLAGS) -o $(BIN_DIR)$(SEP)$(BINARY) ./cmd/api
	@echo Build complete: $(BIN_DIR)$(SEP)$(BINARY)

run:
	go run ./cmd/api
 
clean:
	$(RMDIR) $(BIN_DIR)