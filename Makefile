#!/usr/bin/make -f

# ./build
BUILD_DIR ?= $(CURDIR)/build
NAME := pylons_sdk

export GO111MODULE = on

all: test-all

build: go.sum $(BUILD_DIR)/
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

$(BUILD_DIR)/:
	@mkdir -p $(BUILD_DIR)

go.sum: go.mod
	@go mod verify
	@go mod tidy

clean:
	@rm -rf $(BUILD_DIR)/

.PHONY: build clean

###############################################################################
###                               Commands                                  ###
###############################################################################

reset_chain:
	@echo Resetting chain...
	pylonsd unsafe-reset-all

.PHONY: reset_chain

###############################################################################
###                                Testing                                  ###
###############################################################################

test-all: int_tests fixture_tests

int_tests:
	rm ./cmd/test/nonce.json || true
	go test -v ./cmd/test/ ${ARGS}

fixture_tests:
	rm ./cmd/fixtures_test/nonce.json || true
	go test -v ./cmd/fixtures_test/ ${ARGS}

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@golangci-lint run -c ./.golangci.yml --out-format=tab --issues-exit-code=0

FIND_ARGS := -name '*.go' -type f  -not -path "*.git*"

format:
	@find . $(FIND_ARGS) | xargs gofmt -w -s
	@find . $(FIND_ARGS) | xargs goimports -w -local github.com/Pylons-tech/pylons_sdk

.PHONY: lint format
