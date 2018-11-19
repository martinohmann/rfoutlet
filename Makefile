.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-12s[0m %s\n", $$1, $$2}'

.PHONY: all
all: app binaries ## install dependencies and build everything

.PHONY: app
app: deps-app build-app ## install app dependencies and build

.PHONY: binaries
binaries: deps pack-app build ## install binary dependencies, pack app and build

.PHONY: deps
deps: ## install go deps
	go mod vendor
	go get -u github.com/gobuffalo/packr/packr

.PHONY: deps-app
deps-app: ## install node deps
	cd app && npm install

.PHONY: build
build: ## build binaries
	go build ./cmd/rfoutlet
	go build ./cmd/rfsniff
	go build ./cmd/rftransmit

.PHONY: build-app
build-app: ## build node app
	cd app && yarn build

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
	rm -rf vendor/ app/node_modules/ app/build/
	rm -f rfoutlet rfsniff rftransmit
	packr clean

.PHONY: install
install: ## install go commands into $GOPATH/bin
	go install ./cmd/rfoutlet
	go install ./cmd/rfsniff
	go install ./cmd/rftransmit
