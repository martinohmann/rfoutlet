FROM node:13.1.0-alpine3.10 as node-builder

ADD web/ /web

WORKDIR /web

RUN npm install && \
    npm run build

FROM golang:1.13.4-alpine3.10 as golang-builder

WORKDIR /go/src/github.com/martinohmann/rfoutlet

RUN apk --no-cache add git make

ADD go.mod .
ADD go.sum .
ADD Makefile .

RUN make deps

ARG GOARCH=arm
ARG GOARM=7

COPY --from=node-builder /web/build web/build

ADD cmd/ cmd/
ADD internal/ internal/
ADD pkg/ pkg/
ADD main.go main.go

RUN CGO_ENABLED=0 GOOS=linux make pack-app build

FROM scratch

COPY --from=golang-builder /go/src/github.com/martinohmann/rfoutlet/rfoutlet /rfoutlet

ENTRYPOINT ["/rfoutlet"]
CMD ["serve"]
