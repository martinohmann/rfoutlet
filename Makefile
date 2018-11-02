all: deps deps-app build build-app

.PHONY: deps
deps:
	go get -u github.com/gobuffalo/packr/packr
	go get -u github.com/Masterminds/glide
	glide install

.PHONY: deps-app
deps-app:
	cd app && npm install

.PHONY: build
build:
	packr build ./cmd/rfoutlet
	go build ./cmd/rftransmit

.PHONY: build-app
build-app:
	cd app && yarn build

.PHONY: run
run: build
	./rfoutlet

.PHONY: test
test:
	go test $$(go list ./... | grep -v /vendor/)

.PHONY: coverage
coverage:
	scripts/coverage


.PHONY: clean
clean:
	rm -rf vendor/
	rm rfoutlet rftransmit

.PHONY: install
install:
	packr install ./cmd/rfoutlet
	go install ./cmd/rftransmit
