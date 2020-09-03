SHELL := /bin/bash

lint: ## Go lint the files
	@golint -set_exit_status $$(go list ./...)

fmt: ## Go fmt the files
	@test $$(gofmt -d *.go | tee /dev/stderr | wc -l) -eq 0

vet: ## Go vet the files
	@go vet *.go

test: ## run unit test
	@go test

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
