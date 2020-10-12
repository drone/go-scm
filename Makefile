# Make does not offer a recursive wildcard function, so here's one:
rwildcard=$(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))

GO_DEPENDENCIES := $(call rwildcard,pkg/,*.go) $(call rwildcard,scm/,*.go)
GO := GO111MODULE=on go
GO_NOMOD := GO111MODULE=off go

build: test

test:
	go test ./...

linux: build

.PHONY: check
check: fmt lint sec ## Runs Go format check as well as security checks

get-fmt-deps:
	$(GO_NOMOD) get golang.org/x/tools/cmd/goimports

.PHONY: importfmt
importfmt: get-fmt-deps ## Checks the import format of the Go source files
	@echo "FORMATTING IMPORTS"
	@goimports -w $(GO_DEPENDENCIES)

.PHONY: fmt ## Checks Go source files are formatted properly
fmt: importfmt
	@echo "FORMATTING SOURCE"
	FORMATTED=`$(GO) fmt ./...`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed un-formatted files:\n$(FORMATTED)") || true

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	$(GO_NOMOD) get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(GOLINT) ## Runs 'go vet' anf 'go lint'
	@echo "VETTING"
	$(GO) vet ./...
	@echo "LINTING"
	$(GOLINT) -set_exit_status ./...

GOSEC := $(GOPATH)/bin/gosec
$(GOSEC):
	$(GO_NOMOD) get -u github.com/securego/gosec/cmd/gosec

.PHONY: sec
sec: $(GOSEC) ## Runs gosec to check for potential security issues in the Go source
	@echo "SECURITY SCANNING"
	$(GOSEC) -quiet -fmt=csv ./...
