FROM golang:1.16.2-alpine3.13

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN CGO_ENABLED=0 go mod download
