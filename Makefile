.DEFAULT_GOAL := help

TEST_FLAGS ?= -race
PKGS ?= $(shell go list ./... | grep -v /vendor/)

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-23s[0m %s\n", $$1, $$2}'

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
	cd web && npm run build

.PHONY: pack-app
pack-app: ## pack app using packr
	packr

.PHONY: test
test: ## run tests
	go test $(TEST_FLAGS) $(PKGS)

.PHONY: vet
vet: ## run go vet
	go vet $(PKGS)

.PHONY: coverage
coverage: ## generate code coverage
	go test $(TEST_FLAGS) -covermode=atomic -coverprofile=coverage.txt $(PKGS)
	go tool cover -func=coverage.txt

.PHONY: clean
clean: ## clean dependencies and artifacts
	rm -rf vendor/ web/node_modules/ web/build/
	rm -f rfoutlet
	packr clean

.PHONY: install
install: build ## install rfoutlet into $GOPATH/bin
	mv rfoutlet $(GOPATH)/bin/rfoutlet

.PHONY: images
images: image-amd64 image-armv7 ## build docker images

.PHONY: image-amd64
image-amd64: ## build amd64 image
	docker build --build-arg GOARCH=amd64 -t mohmann/rfoutlet:amd64 .

.PHONY: image-armv7
image-armv7: ## build armv7 image
	docker build --build-arg GOARCH=arm --build-arg GOARM=7 -t mohmann/rfoutlet:armv7 .

.PHONY: gpio-mockup
gpio-mockup: ## create a mock /dev/gpiochip0 using the gpio-mockup kernel module
	sudo modprobe gpio-mockup gpio_mockup_ranges=0,40
