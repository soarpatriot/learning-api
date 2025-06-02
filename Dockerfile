# Start from the official Golang image for building
FROM golang:1.22-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# Install git (required for go mod)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o learning-api main.go

# Start a new minimal image for running
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the built binary from builder
COPY --from=builder /app/learning-api .

# Copy any static files or migrations if needed (optional)
# COPY migrations ./migrations

# Expose the default port
EXPOSE 8080

# Set environment variables for MySQL connection (override as needed)
ENV PROFILE=production

# Run the binary with PROFILE=production
CMD ["/bin/sh", "-c", "PROFILE=production ./learning-api"]
