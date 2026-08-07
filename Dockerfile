FROM golang:1.22-alpine AS builder

WORKDIR /src

ARG TARGETOS=linux
ARG TARGETARCH

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/dach-admin-api \
    ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget \
    && addgroup -S app \
    && adduser -S -D -H -G app app

WORKDIR /app

COPY --from=builder /out/dach-admin-api /app/dach-admin-api

USER app

EXPOSE 8080

ENV PORT=8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${PORT}/health" >/dev/null || exit 1

ENTRYPOINT ["/app/dach-admin-api"]
