all: build test


build: ## Build the application
	@echo "Building..."
	@go build -o bin/app main.go

run: build ## Run the application
	@echo "Running..."
	./bin/app

test: ## Build and run all tests
	@echo "Testing..."
	@go test ./... -v

clean: ## Remove all build artefacts
	@echo "Cleaning..."
	@rm -f main

lint: ## Lint the codebase
	@echo "Linting..."
	@golangci-lint run

format: ## Format the codebase
	@echo "Formatting..."
	@gofmt -w .

help: ## Display available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: help all build run test clean lint format