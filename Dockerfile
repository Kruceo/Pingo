# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pinger main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S pinger && \
    adduser -u 1000 -S pinger -G pinger

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/pinger .

# Copy default config (optional)
COPY --from=builder /app/config.example.json ./config.json

# Change ownership to non-root user
RUN chown -R pinger:pinger /app

# Switch to non-root user
USER pinger

# Expose port (if needed for future web interface)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./pinger"]

# Default command (can be overridden)
CMD ["start", "config.json"]