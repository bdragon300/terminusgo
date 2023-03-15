FROM golang:1.18

WORKDIR /terminusgo
ADD . .
RUN go mod download
