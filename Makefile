.DEFAULT_GOAL := help

GO ?= go
GOLANGCI_LINT_VERSION ?= v1.37.0
TEST_FLAGS ?= -race
IMAGE ?= mohmann/rfoutlet
IMAGE_TAG ?= latest
PKGS ?= $(shell go list ./... | grep -v /vendor/)

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-19s[0m %s\n", $$1, $$2}'

.PHONY: all
all: app binary ## install dependencies and build everything

.PHONY: app
app: deps-app build-app ## install app dependencies and build

.PHONY: binary
binary: deps build ## install binary dependencies and build

.PHONY: deps
deps: ## install go deps
	$(GO) mod download

.PHONY: deps-app
deps-app: ## install node deps
	cd web && npm install

.PHONY: build
build: ## build rfoutlet
	$(GO) build -ldflags="-s -w" -o rfoutlet main.go

.PHONY: build-app
build-app: ## build node app
	cd web && npm run build

.PHONY: test
test: ## run tests
	$(GO) test $(TEST_FLAGS) $(PKGS)

.PHONY: vet
vet: ## run go vet
	$(GO) vet $(PKGS)

.PHONY: coverage
coverage: ## generate code coverage
	$(GO) test $(TEST_FLAGS) -covermode=atomic -coverprofile=coverage.txt $(PKGS)
	$(GO) tool cover -func=coverage.txt

.PHONY: lint
lint: ## run golangci-lint
	command -v golangci-lint > /dev/null 2>&1 || \
	  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	golangci-lint run

.PHONY: clean
clean: ## clean dependencies and artifacts
	rm -rf vendor/ web/node_modules/ web/build/
	rm -f rfoutlet

.PHONY: install
install: build ## install rfoutlet into $GOPATH/bin
	mv rfoutlet $(GOPATH)/bin/rfoutlet

.PHONY: image
image: ## build docker image
	docker build -t $(IMAGE):$(IMAGE_TAG) .

.PHONY: load-gpio-mockup
load-gpio-mockup: ## create a mock /dev/gpiochip0 using the gpio-mockup kernel module
	sudo modprobe gpio-mockup gpio_mockup_ranges=0,40

.PHONY: unload-gpio-mockup
unload-gpio-mockup: ## unload the gpio-mockup kernel module
	sudo modprobe --remove gpio-mockup
