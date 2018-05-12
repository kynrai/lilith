
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

docker-build:
	@cd cmd/$(SERVICE) && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main && \
	docker build -t eu.gcr.io/accuport-prod/$(SERVICE):$(TAG) . && \
	docker tag eu.gcr.io/accuport-prod/$(SERVICE):$(TAG) eu.gcr.io/accuport-prod/$(SERVICE):latest && \
	rm main
.PHONY: docker-build

docker-push:
	@docker push eu.gcr.io/accuport-prod/$(SERVICE):$(TAG)
.PHONY: docker-push

k8s-deploy:
	@helm upgrade accuport-api k8s/charts/accuport --set hash=$(TAG) --install
.PHONY: k8s-deploy

k8s-delete:
	@helm delete accuport-api --purge
.PHONY: k8s-delete

k8s-deploy-dry:
	@helm upgrade accuport-api k8s/charts/accuport --set hash=$(TAG) --dry-run --debug
.PHONY: k8s-deploy

ship: docker-build docker-push k8s-deploy
.PHONY: ship
