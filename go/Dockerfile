# =============================================================================
# Matrix Financial Pipeline - Multi-Stage Docker Build
# =============================================================================
#
# Architecture:
# - Stage 1 (base): Development tools and dependencies
# - Stage 2 (builder): Compile all services
# - Stage 3-N (*-prod): Production images (minimal, one per service)
# - Stage N+1 (dev): Development environment (default)
#
# Usage:
# - Development: docker-compose up (uses 'dev' target)
# - Production:  docker build --target neo-prod -t neo:prod .
# - Testing:     docker run matrix:latest go test ./...
#
# =============================================================================

# =============================================================================
# Stage 1: Base - Development Tools & Dependencies
# =============================================================================
FROM golang:1.24.5 AS base

# Install build tools and code generators
RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/vektra/mockery/v2@latest \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy and download dependencies (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Install testing dependencies
RUN go get github.com/stretchr/testify/assert@v1.11.1 && \
    go get github.com/stretchr/testify/mock@v1.11.1

# =============================================================================
# Stage 2: Builder - Compile All Services
# =============================================================================
FROM base AS builder

# Copy source code
COPY . .

# Build all service binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/neo ./cmd/neo/main.go

# Uncomment as services become available
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/morpheus ./cmd/morpheus/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/trinity ./cmd/trinity/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/oracle ./cmd/oracle/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/tank ./cmd/tank/main.go

# =============================================================================
# Stage 3: Production - Neo Service
# =============================================================================
FROM alpine:latest AS neo-prod

# Add ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the Neo binary
COPY --from=builder /app/bin/neo .

# Expose ports (adjust as needed)
EXPOSE 8080 50051

# Run the service
CMD ["./neo"]

# =============================================================================
# Stage 4: Production - Morpheus Service (Placeholder)
# =============================================================================
FROM alpine:latest AS morpheus-prod

RUN apk --no-cache add ca-certificates
WORKDIR /app

# Uncomment when Morpheus is ready
# COPY --from=builder /app/bin/morpheus .
# EXPOSE 50051
# CMD ["./morpheus"]

CMD ["echo", "Morpheus not yet implemented"]

# =============================================================================
# Stage 5: Production - Trinity Service (Placeholder)
# =============================================================================
FROM alpine:latest AS trinity-prod

RUN apk --no-cache add ca-certificates
WORKDIR /app

# Uncomment when Trinity is ready
# COPY --from=builder /app/bin/trinity .
# EXPOSE 8081
# CMD ["./trinity"]

CMD ["echo", "Trinity not yet implemented"]

# =============================================================================
# Stage 6: Development Environment (Default)
# =============================================================================
FROM base AS dev

# Copy source code (will be overridden by volume mount in docker-compose)
COPY . .

# Build Neo for quick testing
RUN go build -o bin/neo ./cmd/neo/main.go

# Default to bash for interactive development
CMD ["/bin/bash"]

# =============================================================================
# Build Instructions:
# =============================================================================
#
# Development (hot reload with mounted volumes):
#   docker-compose up
#
# Production build for specific service:
#   docker build --target neo-prod -t matrix-neo:latest .
#   docker build --target morpheus-prod -t matrix-morpheus:latest .
#
# Run tests:
#   docker run --rm -v $(pwd):/app -w /app matrix:latest go test ./...
#
# Generate mocks:
#   docker run --rm -v $(pwd):/app -w /app matrix:latest bash cmd/generate-mocks.sh
#
# =============================================================================