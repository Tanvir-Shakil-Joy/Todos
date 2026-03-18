# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

# Build the binary with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# ---- Run Stage ----
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget && \
    addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
