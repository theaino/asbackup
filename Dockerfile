FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o asbackup .

CMD ["/build/asbackup"]
