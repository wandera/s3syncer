# Builder image
FROM golang:1.20.6-alpine3.17 AS builder

WORKDIR /build

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build GOMODCACHE=/go/pkg/mod GOCACHE=/root/.cache/go-build go build

# Runtime image
FROM alpine:3.18.2
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /build/s3syncer /app/s3syncer
WORKDIR /app

ENTRYPOINT ["./s3syncer"]
