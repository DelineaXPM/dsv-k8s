FROM alpine:latest AS build
RUN apk update && apk upgrade
RUN apk add go

WORKDIR /b
COPY cmd/ ./cmd
COPY pkg/ ./pkg
COPY go.mod go.sum ./
RUN go build cmd/injector-svc.go

ARG cert_file
ARG key_file
ARG roles_file

COPY ${cert_file} ./dsv.pem
COPY ${key_file} ./dsv.key
COPY ${roles_file} ./roles.json

FROM alpine:latest
RUN apk update && apk upgrade
RUN addgroup dsv && adduser -S -G dsv dsv

WORKDIR /home/dsv
COPY --chown=dsv:dsv --from=build /b/injector-svc /b/dsv.pem /b/dsv.key /b/roles.json /home/dsv/

USER dsv

ENTRYPOINT ["./injector-svc"]
