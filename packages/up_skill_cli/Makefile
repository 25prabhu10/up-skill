DBG_MAKEFILE ?=
ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
else
    # If we're not debugging the Makefile, don't echo recipes.
    MAKEFLAGS += -s
endif

# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
# Be pedantic about undefined variables.
MAKEFLAGS += --warn-undefined-variables

# Project variables
BINARY_NAME := up-skill
SRCDIR := .
BINDIR := bin

# VCS / build metadata
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
RELEASE ?= 0

# Go commands
GO ?= go
GO_BUILD := CGO_ENABLED=0 $(GO) build -trimpath
GO_LINT := $(shell command -v golangci-lint 2> /dev/null)

# Linker flags
PKG := "$(shell $(GO) list -m)/pkg/build_info"
LDFLAGS := -X $(PKG).AppName=$(BINARY_NAME) \
	   -X $(PKG).Version=$(VERSION) \
	   -X $(PKG).GitCommit=$(GIT_COMMIT) \
	   -X $(PKG).BuildDate=$(BUILD_DATE)

ifeq ($(RELEASE),1)
	LDFLAGS += -s -w
endif

.DEFAULT_GOAL := all

# check if `go` is installed
ifeq (, $(shell which $(GO)))
$(error "No 'go' binary found in PATH. Please install Go to proceed.")
endif

# Delete target file if the recipe fails
.DELETE_ON_ERROR:

.PHONY: all
all: build test

# Create bin directory if it doesn't exist
$(BINDIR):
	@mkdir -p $(BINDIR)

.PHONY: build
build: $(BINDIR) ## Build the binary for current OS/Arch
	$(GO_BUILD) -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(BINARY_NAME) $(SRCDIR)

.PHONY: release
release: ## Build release binaries for multiple platforms (linux, darwin, windows)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			echo "Building $$os-$$arch..."; \
			extension=""; [ "$$os" = "windows" ] && extension=".exe"; \
			GOOS=$$os GOARCH=$$arch $(GO_BUILD) -ldflags "$(LDFLAGS)" \
				-o $(BINDIR)/$(BINARY_NAME).$$os.$$arch.v$(VERSION)$$extension $(SRCDIR); \
		done; \
	done

.PHONY: run
run: ## Run the application
	$(GO) run -ldflags "$(LDFLAGS)" $(SRCDIR) $(ARGS)

.PHONY: test
test: ## Run tests
	$(GO) test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: format
fmt: ## Format files
	$(GO) fmt ./...

.PHONY: lint
lint: ## Run vet and linter (requires golangci-lint)
	$(GO) vet ./...
	@if [ -z "$(GO_LINT)" ]; then \
		error "golangci-lint not installed, skipping..."; \
	else \
		$(GO_LINT) run ./...; \
	fi

.PHONY: tidy
tidy: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: update-deps
update-deps: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: clean
clean: ## Remove build artifacts
	@echo "Removing build artifacts..."
	@rm -rf $(BINDIR) coverage.out coverage.html

.PHONY: install
install: build ## Install the binary to $GOPATH/bin
	@(GO) install -ldflags "$(LDFLAGS)" $(SRCDIR)
	@echo "Installed to $(shell $(GO) env GOPATH)/bin/$(BINARY_NAME)"

.PHONY: uninstall
uninstall: ## Remove installed binary from GOBIN/GOPATH
	@rm -f $(shell $(GO) env GOPATH)/bin/$(BINARY_NAME)

.PHONY: help
help: ## Display this help message
	@echo "up-skill CLI - Makefile commands"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)