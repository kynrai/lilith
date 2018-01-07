
BINARIES=$$(go list ./cmd/...)
TESTABLE=$$(go list ./... | grep -v /vendor/)

all: vet test build clean

deps:
	@dep ensure && dep ensure -update
.PHONY: deps

build:
	@go install -v  $(BINARIES)
.PHONY: build

test:
	@go test -v $(TESTABLE)
.PHONY: test

vet:
	@go vet $(TESTABLE)
.PHONY: vet

fmt:
	@goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
.PHONY: fmt

lint:
	@golint $(TESTABLE)
.PHONY: lint

clean:
	@go clean
.PHONY: clean

local: all
	@heroku local web

