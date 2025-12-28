# =========================
# Stage 1: Builder
# =========================
FROM golang:1.24.4-alpine AS builder

# Install only required build deps
RUN apk add --no-cache \
    ca-certificates \
    build-base

WORKDIR /app

# Enable modules explicitly
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy go mod files first (better cache)
COPY go.mod go.sum ./

# Download deps (vendor optional but respected)
RUN go mod download

# Copy source code
COPY vendor ./vendor
COPY account ./account
COPY catalog ./catalog
COPY order ./order

# Build binary
RUN go build \
    -mod=vendor \
    -trimpath \
    -ldflags="-s -w" \
    -o app \
    ./order/cmd/order

# =========================
# Stage 2: Runtime
# =========================
FROM alpine:3.20

# Add certs + non-root user
RUN apk add --no-cache ca-certificates \
    && adduser -D -g '' appuser

WORKDIR /app

# Copy binary
COPY --from=builder /app/app .

# Switch to non-root
USER appuser

EXPOSE 8080

ENTRYPOINT ["./app"]
