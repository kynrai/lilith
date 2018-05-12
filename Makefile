
BINARIES=$$(go list ./cmd/...)
TESTABLE=$$(go list ./...)

all: test build

deps:
	@dep ensure && dep ensure -update
.PHONY: deps

build:
	@go install -v  $(BINARIES)
.PHONY: build

test:
	@go test -v $(TESTABLE)
.PHONY: test

fmt:
	@goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
.PHONY: fmt

lint:
	@golint $(TESTABLE)
.PHONY: lint

local: all
	@heroku local web

