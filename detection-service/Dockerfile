# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY cmd  ./cmd/
COPY config  ./config/
COPY internal  ./internal/
COPY helpers ./helpers/
RUN go mod download

RUN go build -o app ./cmd/detection-service/detection-service.go
CMD ["./app"]

EXPOSE 10001