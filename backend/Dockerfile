# Multi-stage build for Go backend with Python parser

# Stage 1: Build Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Final image with Python support
FROM python:3.11-slim

WORKDIR /app

# Install Python dependencies
RUN pip install --no-cache-dir biopython==1.81

# Copy Go binary from builder
COPY --from=builder /app/main .

# Copy Python scripts
COPY scripts/ ./scripts/

# Create upload directory
RUN mkdir -p uploads

# Copy .env file (optional)
COPY .env* ./

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]