# Build stage: Go 1.25 on Alpine
FROM golang:1.25-alpine AS builder

# Install git if your modules need it
# RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o proxy-server .

# Final stage: small runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/proxy-server .

# Expose port
EXPOSE 6443

# Run the binary
CMD ["./proxy-server", "-host", "127.0.0.1", "-port", "6443"]