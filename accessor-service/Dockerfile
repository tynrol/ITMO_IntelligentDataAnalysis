# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /build
ENV UNSPLASH_TOKEN="L5Z1WbsdoD7teDNFsukmJ3k32rGkDaphX4c3DTNdb3I"
COPY go.mod .
COPY go.sum .
COPY cmd  ./cmd/
COPY config  ./config/
COPY init/ ./init/
COPY internal  ./internal/
RUN go mod download

RUN go build -o app ./cmd/accessor-service/accessor-service.go
CMD ["./app"]

EXPOSE 10000