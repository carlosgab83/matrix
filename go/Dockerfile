# =============================================================================
# Matrix Financial Pipeline - Docker Build Environment
# =============================================================================
#
# ⚠️  WORK IN PROGRESS - This Dockerfile is incomplete
#
# Current Status:
# ✅ Base Go environment configured
# ✅ Protocol Buffers and gRPC tools installed
# ✅ Neo service build configuration
# 🚧 Trinity build configuration - TODO
# ⏳ Morpheus build configuration - Planned
# ⏳ Oracle build configuration - Planned
# ⏳ Tank build configuration - Planned
# ⏳ Multi-stage builds for production - Planned
# ⏳ Service-specific optimizations - Planned
#
# =============================================================================

# Base image with Go 1.24.5
FROM golang:1.24.5

# Install protoc, gRPC plugins, and Mockery
RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/vektra/mockery/v2@latest

# Set working directory
WORKDIR /app

# Copy module files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies (caches this layer)
RUN go mod download

RUN go get github.com/stretchr/testify/assert@v1.11.1 && \
    go get github.com/stretchr/testify/mock@v1.11.1

# Copy the rest of the source code
COPY . .

# =============================================================================
# BUILD BINARIES
# =============================================================================

# ✅ Neo - Stock Price Collector (Available)
RUN go build -o bin/neo-app ./cmd/neo/main.go

# 🚧 Trinity - News Collector (Basic implementation)
# TODO: Uncomment when ready
# RUN go build -o bin/trinity-app ./cmd/trinity/main.go

# ⏳ Morpheus - Data Hub (Not implemented yet)
# TODO: Add when service is ready
# RUN go build -o bin/morpheus-app ./cmd/morpheus/main.go

# ⏳ Oracle - Analytics Engine (Not implemented yet)
# TODO: Add when service is ready
# RUN go build -o bin/oracle-app ./cmd/oracle/main.go

# ⏳ Tank - Notification Gateway (Not implemented yet)
# TODO: Add when service is ready
# RUN go build -o bin/tank-app ./cmd/tank/main.go

# =============================================================================
# RUNTIME CONFIGURATION
# =============================================================================

# Expose ports for services
# EXPOSE 8080  # Neo HTTP API (if implemented)
# EXPOSE 50051 # gRPC port (Morpheus/Trinity)
# EXPOSE 9090  # Metrics port (if implemented)

# Default command - Currently only Neo is available
# TODO: Add service selection mechanism or docker-compose setup
# CMD ["./bin/neo-app"]

# For development: Start with bash for manual testing
CMD ["/bin/bash"]