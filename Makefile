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
.PHONY: bootstrap
bootstrap: ## Downloads and cleans up all dependencies
	@go mod tidy

.PHONY: fmt
fmt: ## Formats go files
	@gofmt -w -s $(GO_FILES)

.PHONY: check
check: ## Checks code for linting/construct errors
	@gofmt -l $(GO_FILES)
	@go list -f '{{.Dir}}' ./... | xargs go vet;

.PHONY: test
test: check ## Runs all tests
	@go test --race $(TEST) -parallel=20

.PHONY: benchmark
benchmark: check ## Runs all benchmarks
	@go test -bench=. -benchmem

.PHONY: coverage
coverage: test ## Runs the tests and shows the code coverage report
	@mkdir -p .cover
	@go test $(TEST) -race -coverprofile=.cover/cover.out -covermode=atomic
	@go tool cover -html=.cover/cover.out

.PHONY: clean
clean: ## Cleans up temporary and compiled files
	@rm -rf .cover

.PHONY: help
help: ## Shows available targets
	@fgrep -h "## " $(MAKEFILE_LIST) | fgrep -v fgrep | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-13s\033[0m %s\n", $$1, $$2}'
