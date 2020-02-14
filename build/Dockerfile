FROM alpine:latest AS build
RUN apk update && apk upgrade
RUN apk add go

WORKDIR /b
COPY cmd/ ./cmd
COPY pkg/ ./pkg
COPY go.mod go.sum ./
RUN go build cmd/injector-svc.go

FROM alpine:latest
RUN apk update && apk upgrade
RUN addgroup dsv && adduser -S -G dsv dsv

WORKDIR /home/dsv
COPY --chown=dsv:dsv --from=build /b/injector-svc .

USER dsv
