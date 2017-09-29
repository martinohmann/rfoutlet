PROGRAM_NAME := rf-outlet

DEPLOY_ROOT ?= /tmp/$(PROGRAM_NAME)

default: build

.PHONY: deps test build

test:
	go test $$(go list ./... | grep -v /vendor/)

coverage:
	scripts/coverage

deps:
	go get github.com/Masterminds/glide
	glide install

build:
	go build

run:
	go run main.go

clean:
	rm -rf vendor/
	rm $(PROGRAM_NAME)
