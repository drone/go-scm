.PHONY: test

PACKAGES := $(shell go list ./... | grep -v '/integration')

test:
	go test -v -coverprofile=coverage.out -covermode=atomic $(PACKAGES)
