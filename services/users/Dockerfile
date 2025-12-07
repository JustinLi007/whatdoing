#syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:latest AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
RUN go build -ldflags="-s -w" -o /usr/bin/app ./cmd/app/main.go

FROM alpine:latest AS prod
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /usr/bin/app /usr/bin/app
ENTRYPOINT ["/usr/bin/app"]
