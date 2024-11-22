# Makefile for Go project

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

# Targets
all: test build

build:
	$(GOBUILD) -o bin/main cmd/main.go

run:
	$(GORUN) cmd/main.go

test:
	$(GOTEST) -v ./...

testrun:
	$(GOTEST) -v ./...
	$(GORUN) cmd/main.go

clean:
	$(GOCLEAN)
	rm -rf bin/

.PHONY: all build test clean
