# Builder stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Action binary
RUN go build -o /action cmd/action/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /action /app/action
# Copy the default configuration file
COPY --from=builder /app/analizer_config.yaml /app/analizer_config.yaml

# Set the entrypoint
ENTRYPOINT ["/app/action"]
