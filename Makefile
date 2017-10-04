PROGRAM_NAME := rf-outlet

DEPLOY_ROOT ?= /tmp/$(PROGRAM_NAME)

default: build

.PHONY: test
test:
	go test $$(go list ./... | grep -v /vendor/)

.PHONY: coverage
coverage:
	scripts/coverage

.PHONY: deps
deps:
	go get github.com/Masterminds/glide
	glide install
	cd frontend && npm install

.PHONY: build
build:
	go build
	cd frontend && yarn build

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	rm -rf vendor/
	rm $(PROGRAM_NAME)
