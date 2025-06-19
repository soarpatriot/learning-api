# Start from the official Golang image for building
FROM golang:1.24-alpine AS builder

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
RUN go build -o main main.go

# Start a new minimal image for running
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates
COPY RapidSSLTLSRSACAG1.crt /usr/local/share/ca-certificates/RapidSSLTLSRSACAG1.crt
RUN update-ca-certificates
# Expose the default port
EXPOSE 8000


WORKDIR /opt/application

# Copy the built binary from builder
COPY --from=builder /app/main /opt/application/main
COPY --from=builder /app/run.sh /opt/application/run.sh
COPY --from=builder /app/config.yaml /opt/application/config.yaml

# Copy any static files or migrations if needed (optional)
# COPY migrations ./migrations

USER root
RUN chmod -R 777 /opt/application/run.sh

CMD /opt/application/run.sh

