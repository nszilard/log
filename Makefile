#!make
#----------------------------------------
# Settings
#----------------------------------------
.DEFAULT_GOAL := help

#--------------------------------------------------
# Variables
#--------------------------------------------------
GO_FILES?=$$(find . -name '*.go')
TEST?=$$(go list ./... | grep -v /vendor/)

#--------------------------------------------------
# Targets
#--------------------------------------------------
bootstrap: ## Downloads and cleans up all dependencies
	@go mod tidy

fmt: ## Formats go files
	@echo "==> Formatting files..."
	@gofmt -w -s $(GO_FILES)
	@echo ""

check: ## Checks code for linting/construct errors
	@echo "==> Checking if files are well formatted..."
	@gofmt -l $(GO_FILES)
	@echo ""
	@echo "==> Checking if files pass go vet..."
	@go list -f '{{.Dir}}' ./... | xargs go vet;
	@echo ""

test: check ## Runs all tests
	@echo "==> Running tests..."
	@go test --race $(TEST) -parallel=20
	@echo ""

coverage: ## Runs code coverage
	@mkdir -p .cover
	@go test $(TEST) -race -coverprofile=.cover/cover.out -covermode=atomic

show-coverage: coverage ## Shows code coverage report in your web browser
	@go tool cover -html=.cover/cover.out

.PHONY: bootstrap check package fmt test coverage show-coverage clean help

clean: ## Cleans up temporary and compiled files
	@echo "==> Cleaning up ..."
	@rm -rf .cover
	@echo "    [✓]"
	@echo ""

help: ## Shows available targets
	@fgrep -h "## " $(MAKEFILE_LIST) | fgrep -v fgrep | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-13s\033[0m %s\n", $$1, $$2}'
