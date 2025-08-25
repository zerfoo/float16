SHELL := /bin/sh

PKG := ./...
COVER_DIR := coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML := $(COVER_DIR)/coverage.html

.PHONY: all test race vet fmt fmt-check lint lint-fix cover cover-html ci

all: test

# Run unit tests
test:
	go test $(PKG)

# Run tests with race detector
race:
	go test -race $(PKG)

# Static analysis
vet:
	go vet $(PKG)

# Auto-format code
fmt:
	gofmt -s -w .

# Check formatting without modifying files; fails if formatting needed
fmt-check:
	@diff=$$(gofmt -s -l .); \
	if [ -n "$$diff" ]; then \
		echo "Files need formatting:"; echo "$$diff"; exit 1; \
	else \
		echo "Formatting OK"; \
	fi

# Lint: go vet + formatting check + optional golangci-lint if installed
lint:
	@echo "Running go vet"; go vet $(PKG)
	@echo "Checking formatting"; \
	diff=$$(gofmt -s -l .); if [ -n "$$diff" ]; then echo "Files need formatting:"; echo "$$diff"; exit 1; else echo "Formatting OK"; fi
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint"; golangci-lint run; \
	else \
		echo "golangci-lint not installed; skipping"; \
	fi

# Attempt to fix lint: gofmt + optional golangci-lint --fix if installed
lint-fix: fmt
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint --fix"; golangci-lint run --fix || true; \
	else \
		echo "golangci-lint not installed; skipping"; \
	fi

# Generate coverage profile and print total coverage
cover:
	mkdir -p $(COVER_DIR)
	go test -covermode=atomic -coverprofile=$(COVER_PROFILE) $(PKG)
	go tool cover -func=$(COVER_PROFILE) | tail -n 1

# Generate HTML coverage report
cover-html: cover
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Wrote $(COVER_HTML)"

# CI-style aggregate target
ci: fmt-check vet test cover
