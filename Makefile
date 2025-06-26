# Makefile for Simserver

# Get build information
GIT_HASH := $(shell git show -s --format=%H)
BUILD_TIME := $(shell git show -s --format=%cd)
GO_VERSION := $(shell go version)

# Set ldflags to embed build information
LDFLAGS := -ldflags "-s -w -X 'main.gitHash=$(GIT_HASH)' -X 'main.buildTime=$(BUILD_TIME)' -X 'main.goVersion=$(GO_VERSION)'"

# Set the output binary name
BINARY_NAME := Simserver
ifeq ($(OS),Windows_NT)
    BINARY_NAME := Simserver.exe
endif

.PHONY: all build package clean

all: build

build:
	@echo "Building Simserver..."
	@go build $(LDFLAGS) -o $(BINARY_NAME)

package: build
	@echo "Packaging Simserver..."
	@mkdir -p dist
	@tar -czvf dist/Simserver-$(shell uname -s)-$(shell uname -m).tar.gz $(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist
