FROM node:10-alpine as node-builder

ADD app/ /app

WORKDIR /app

RUN npm install && \
    yarn build

FROM golang:1.9-alpine as golang-builder

WORKDIR /go/src/github.com/martinohmann/rfoutlet

ADD glide.lock .
ADD glide.yaml .

RUN apk --no-cache add git && \
    go get -u github.com/gobuffalo/packr/packr && \
	go get -u github.com/Masterminds/glide && \
	glide install

COPY --from=node-builder /app/build app/build

ADD cmd/ cmd/
ADD internal/ internal/
ADD pkg/ pkg/

ARG GOARCH=arm
ARG GOARM=7

RUN CGO_ENABLED=0 GOOS=linux packr build ./cmd/rfoutlet

FROM scratch

COPY --from=golang-builder /go/src/github.com/martinohmann/rfoutlet/rfoutlet /rfoutlet

ENTRYPOINT ["/rfoutlet"]
