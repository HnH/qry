.PHONY: help prepare lint test coverage clean
.DEFAULT_GOAL := help

prepare: ## Prepare build prerequisities
	GO111MODULE=on go get

clean: ## clean working tree
	rm qry coverage.out

lint: ## Run the linters
	go vet ./...
	revive -formatter friendly -config .circleci/revive.toml ./...
	gosec -tests -vendor ./...

test: ## Run tests
	go test -race -cover -coverprofile=coverage.out ./...

coverage: ## Show code coverage with tests
	go tool cover -html=coverage.out

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
