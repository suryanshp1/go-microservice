# ============================
# Stage 1: Builder
# ============================
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    git

WORKDIR /app

# Enable Go modules explicitly
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy only go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy source code
COPY account ./account

# Build binary
RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o app \
    ./account/cmd/account

# ============================
# Stage 2: Runtime
# ============================
FROM alpine:3.20

# Security + certificates
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary
COPY --from=builder /app/app .

# Drop privileges
USER appuser

EXPOSE 8080

ENTRYPOINT ["./app"]
