.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-12s[0m %s\n", $$1, $$2}'

.PHONY: all
all: deps deps-app build-app build ## install dependencies and build everything

.PHONY: deps
deps: ## install go deps
	go get -u github.com/gobuffalo/packr/packr
	go get -u github.com/Masterminds/glide
	glide install

.PHONY: deps-app
deps-app: ## install node deps
	cd app && npm install

.PHONY: build
build: ## build go commands
	packr build ./cmd/rfoutlet
	go build ./cmd/rfsniff
	go build ./cmd/rftransmit

.PHONY: build-app
build-app: ## build node app
	cd app && yarn build

.PHONY: test
test: ## run tests
	go test $$(go list ./... | grep -v /vendor/)

.PHONY: coverage
coverage: ## generate code coverage
	scripts/coverage

.PHONY: clean
clean: ## clean dependencies and artifacts
	rm -rf vendor/ app/node_modules/ app/build/
	rm rfoutlet rfsniff rftransmit

.PHONY: install
install: ## install go commands into $GOPATH/bin
	packr install ./cmd/rfoutlet
	go install ./cmd/rfsniff
	go install ./cmd/rftransmit

.PHONY: images
images: images-amd64 images-armv7 ## build docker images

.PHONY: image-amd64
images-amd64: ## build amd64 images
	docker build --build-arg GOARCH=amd64 -t mohmann/rfoutlet:amd64 .

.PHONY: image-armv7
images-armv7: image-rfoutlet-armv7 image-rftransmit-armv7 ## build armv7 images

.PHONY: image-rfoutlet-armv7
image-rfoutlet-armv7: ## build rfoutlet armv7 image
	docker build -t mohmann/rfoutlet:armv7 .

.PHONY: image-rftransmit-armv7
image-rftransmit-armv7: ## build rftransmit armv7 image
	docker build -t mohmann/rftransmit:armv7 -f Dockerfile.rftransmit .
