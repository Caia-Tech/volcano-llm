# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o volcano-llm ./cmd/temporal-server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' volcano

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/volcano-llm .

# Copy configuration files
COPY --chown=volcano:volcano ./repos /app/repos

# Switch to non-root user
USER volcano

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./volcano-llm"]