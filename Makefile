.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-23s[0m %s\n", $$1, $$2}'

.PHONY: all
all: app binaries ## install dependencies and build everything

.PHONY: app
app: deps-app build-app ## install app dependencies and build

.PHONY: binaries
binaries: deps pack-app build ## install binary dependencies, pack app and build

.PHONY: deps
deps: ## install go deps
	go mod download
	go get github.com/gobuffalo/packr/packr@v1.30.1

.PHONY: deps-app
deps-app: ## install node deps
	cd web && npm install

.PHONY: build
build: ## build binaries
	go build ./cmd/rfoutlet
	go build ./cmd/rfsniff
	go build ./cmd/rftransmit

.PHONY: build-app
build-app: ## build node app
	cd web && yarn build

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
	rm -f rfoutlet rfsniff rftransmit
	packr clean

.PHONY: install
install: ## install go commands into $GOPATH/bin
	go install ./cmd/rfoutlet
	go install ./cmd/rfsniff
	go install ./cmd/rftransmit

.PHONY: images
images: images-amd64 images-armv7 ## build docker images

.PHONY: image-amd64
images-amd64: ## build amd64 images
	docker build --build-arg GOARCH=amd64 -t mohmann/rfoutlet:amd64 .

.PHONY: image-armv7
images-armv7: image-rfoutlet-armv7 image-rfsniff-armv7 image-rftransmit-armv7 ## build armv7 images

.PHONY: image-rfoutlet-armv7
image-rfoutlet-armv7: ## build rfoutlet armv7 image
	docker build -t mohmann/rfoutlet:armv7 .

.PHONY: image-rfsniff-armv7
image-rfsniff-armv7: ## build rfsniff armv7 image
	docker build -t mohmann/rfsniff:armv7 -f Dockerfile.rfsniff .

.PHONY: image-rftransmit-armv7
image-rftransmit-armv7: ## build rftransmit armv7 image
	docker build -t mohmann/rftransmit:armv7 -f Dockerfile.rftransmit .
