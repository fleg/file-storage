export GOBIN ?= $(shell pwd)/bin

GODOTENV = $(GOBIN)/godotenv
REFLEX = $(GOBIN)/reflex

.PHONY: install
install:
	echo $$GOBIN
	go mod download
	go install github.com/joho/godotenv/cmd/godotenv@v1.5.0
	go install github.com/cespare/reflex@v0.3.1

.PHONY: sync
sync:
	go mod tidy

.PHONY: dev
dev:
	$(GODOTENV) -f .env $(REFLEX) -d none -r \.go$$ -s go run cmd/app.go

.PHONY: start
start:
	go run cmd/app.go

.PHONY: fmt
fmt:
	gofmt -s -w -e .
