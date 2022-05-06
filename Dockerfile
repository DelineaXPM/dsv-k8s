FROM alpine:latest AS build
RUN apk update && apk upgrade
RUN apk add go

WORKDIR /build
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg
COPY go.mod go.sum ./
RUN go build -o dsv-injector ./cmd/injector
RUN go build -o dsv-syncer ./cmd/syncer

FROM alpine:latest
RUN apk update && apk upgrade
RUN addgroup dsv && adduser -S -G dsv dsv

COPY --from=build /build/dsv-injector /build/dsv-syncer /usr/bin/

USER dsv
WORKDIR /home/dsv
