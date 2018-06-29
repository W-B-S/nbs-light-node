# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_BASE=nbs-light-node
VERSION = 0.01
GOFILES=$(wildcard *.go)

OS ?= $(shell sh -c 'uname -s 2>/dev/null || echo not')
ifeq ($(OS),Windows_NT)
	BINARY_NAME = $(BINARY_BASE)_win
else
	BINARY_NAME = $(BINARY_BASE)_mac
endif

dir := core
include $(dir)/Rules.mk

all: test build
.PHONY: all

build: vet
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
.PHONY: build

test:
	$(GOTEST) -v ./...
.PHONY: test

run: build
	./$(BINARY_NAME)
.PHONY: run

#deps:
#	$(GOGET) github.com/markbates/goth
#.PHONY: deps

vet:
	GOVET ./...

clean:
	$(GOCLEAN)
	if [ -f $(BINARY_NAME) ] ; then rm -rf $(BINARY_NAME); fi
.PHONY: clean