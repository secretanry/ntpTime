BINARY_NAME=ntpTime

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOLINT=golint
GOLINT_PATH=$(shell go env GOPATH)/bin/golint

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "NTP Time Client - Testing Commands:"
	@echo ""
	@echo "make test          - Run all tests"
	@echo "make test-coverage - Run tests with coverage"
	@echo "make vet          - Run go vet"
	@echo "make lint         - Run golint"
	@echo "make check        - Run all quality checks"
	@echo "make build        - Build the application"
	@echo "make run          - Run the application"
	@echo "make clean        - Clean build artifacts"

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) .

.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run .

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v .

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -cover .

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) .

.PHONY: lint
lint:
	@echo "Installing golint if not present..."
	@if [ ! -f "$(GOLINT_PATH)" ]; then \
		echo "Installing golint..."; \
		$(GOCMD) install golang.org/x/lint/golint@latest; \
	fi
	@echo "Running golint..."
	$(GOLINT_PATH) .

.PHONY: check
check: vet lint
	@echo "All quality checks passed!"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -f $(BINARY_NAME)

