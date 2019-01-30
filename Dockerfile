# Builder image
FROM golang:1.11 AS builder

WORKDIR /build

# Docker Cloud args, from hooks/build.
ARG CACHE_TAG
ENV CACHE_TAG ${CACHE_TAG}

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

# Runtime image
FROM alpine:3.8
COPY --from=builder /build/s3syncer /app/s3syncer
WORKDIR /app

ENTRYPOINT ["./s3syncer"]