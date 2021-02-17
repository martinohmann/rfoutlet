FROM node:15.2.0-alpine3.10 as node-builder

ADD web/ /web

WORKDIR /web

RUN npm install && \
    npm run build

FROM golang:1.16.0-alpine3.12 as golang-builder

WORKDIR /go/src/github.com/martinohmann/rfoutlet

RUN apk --no-cache add git make

ADD go.mod .
ADD go.sum .
ADD Makefile .

RUN make deps

ARG GOARCH
ARG GOARM

COPY --from=node-builder /web/build web/build

ADD cmd/ cmd/
ADD internal/ internal/
ADD pkg/ pkg/
ADD web/static.go web/static.go
ADD main.go main.go

RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch

COPY --from=golang-builder /go/src/github.com/martinohmann/rfoutlet/rfoutlet /rfoutlet

ENTRYPOINT ["/rfoutlet"]
CMD ["serve"]
