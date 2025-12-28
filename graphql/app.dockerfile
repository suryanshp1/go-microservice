# =========================================================
# Stage 1: Build
# =========================================================
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    ca-certificates \
    git \
    build-base

WORKDIR /app

# Cache dependencies first (best practice)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY account ./account
COPY catalog ./catalog
COPY order ./order
COPY graphql ./graphql

# Build a statically linked binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o app ./graphql

# =========================================================
# Stage 2: Runtime
# =========================================================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary only
COPY --from=builder /app/app ./app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./app"]
