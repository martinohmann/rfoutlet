.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-12s[0m %s\n", $$1, $$2}'

.PHONY: all
all: app binary ## install dependencies and build everything

.PHONY: app
app: deps-app build-app ## install app dependencies and build

.PHONY: binary
binary: deps pack-app build ## install binary dependencies, pack app and build

.PHONY: deps
deps: ## install go deps
	go mod download
	go get github.com/gobuffalo/packr/packr@v1.30.1

.PHONY: deps-app
deps-app: ## install node deps
	cd web && npm install

.PHONY: build
build: ## build rfoutlet
	go build -ldflags="-s -w" -o rfoutlet main.go

.PHONY: build-app
build-app: ## build node app
	cd web && npm build

.PHONY: pack-app
pack-app: ## pack app using packr
	packr

.PHONY: test
test: ## run tests
	go test -race -tags="$(TAGS)" $$(go list ./... | grep -v /vendor/)

.PHONY: vet
vet: ## run go vet
	go vet $$(go list ./... | grep -v /vendor/)

.PHONY: coverage
coverage: ## generate code coverage
	scripts/coverage

.PHONY: clean
clean: ## clean dependencies and artifacts
	rm -rf vendor/ web/node_modules/ web/build/
	rm -f rfoutlet
	packr clean

.PHONY: install
install: ## install go commands into $GOPATH/bin
	go install -ldflags="-s -w" main.go
