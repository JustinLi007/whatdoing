#syntax=docker/dockerfile:1

FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /usr/bin/app ./cmd/app/main.go

FROM alpine:latest AS prod
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /usr/bin/app /usr/bin/app
ENTRYPOINT ["/usr/bin/app"]
