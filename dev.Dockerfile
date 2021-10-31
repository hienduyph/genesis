FROM golang:1.17-buster

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download
