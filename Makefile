.PHONY: test

test:
	go test -v -coverprofile=coverage.out -covermode=atomic ./... || true
