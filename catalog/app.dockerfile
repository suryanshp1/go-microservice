# ------------------------------------------------------------
# Stage 1: Build
# ------------------------------------------------------------
FROM golang:1.24.4-alpine AS builder

# Install only what is required for build
RUN apk add --no-cache ca-certificates git

# Enable reproducible builds
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy only dependency files first (better cache)
COPY go.mod go.sum ./

# Download deps separately for caching
RUN go mod download

# Copy source code
COPY catalog ./catalog

# Build binary
RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o app \
    ./catalog/cmd/catalog

# ------------------------------------------------------------
# Stage 2: Runtime
# ------------------------------------------------------------
FROM gcr.io/distroless/base-debian12 AS runtime

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Use non-root user (security)
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/app"]
