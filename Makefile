.DEFAULT_GOAL := help

LOCAL_BIN=$(CURDIR)/

GOENV:=GOPRIVATE=""

.PHONY: run
run: ## run project
	$(GOENV) go run main.go --config config.json

.PHONY: test
test: ## run tests in project
	$(GOENV) go test ./...

.PHONY: release-mac
release-mac: ## build .dmg app
	fyne package -os darwin -name LMS  -icon ./data/app-icon.png

.PHONY: download
download: ## dowmload deps
	$(GOENV) go mod download

.PHONY: tidy
tidy: ## check deps
	$(GOENV) go mod tidy

.PHONY: help
help:
	@grep --no-filename -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'