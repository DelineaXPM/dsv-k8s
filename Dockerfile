FROM alpine:latest AS build
RUN apk update && apk upgrade
RUN apk add go

WORKDIR /b
COPY cmd/ ./cmd
COPY pkg/ ./pkg
COPY go.mod go.sum ./
RUN go build cmd/dsv-injector-svc.go

FROM alpine:latest
RUN apk update && apk upgrade
RUN addgroup dsv && adduser -S -G dsv dsv

COPY --from=build /b/dsv-injector-svc /usr/bin/

USER dsv
WORKDIR /home/dsv
