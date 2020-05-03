SHELL := /bin/bash

all: install

go.sum: go.mod
	GO111MODULE=on go mod verify

init_accounts:
	bash ./init_accounts.sh

reset_chain:
	pylonsd unsafe-reset-all

int_tests:
	rm ./cmd/test/nonce.json || true
	go test -v ./cmd/test/ ${ARGS}

fixture_tests:
	rm ./cmd/fixtures_test/nonce.json || true
	go test -v ./cmd/fixtures_test/ ${ARGS}
