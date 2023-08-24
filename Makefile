# Makefile for cross-compiling Go code for RPi3

PROGRAM_NAME := doorbell

# RPi cross-compile settings
GOOS := linux
GOARCH := arm
GOARM := 7

# RPi SSH connection details (modify these as necessary)
PI_USER := pi
PI_HOST := doorbell.local
PI_DIRECTORY := /home/pi/

all: build

# Cross-compilation for RPi3
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o $(PROGRAM_NAME) $(PROGRAM_NAME).go

# Transfer the binary to the RPi
transfer: build
	scp $(PROGRAM_NAME) $(PI_USER)@$(PI_HOST):$(PI_DIRECTORY)

clean:
	rm -f $(PROGRAM_NAME)

.PHONY: all build transfer clean
