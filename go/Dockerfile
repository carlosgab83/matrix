# Base image with Go 1.24.5
FROM golang:1.24.5

# Install protoc and gRPC plugins
RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set working directory
WORKDIR /app

# Copy module files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies (caches this layer)
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build binaries
RUN go build -o bin/neo-app ./cmd/neo/main.go

# Run the binary by default
# CMD ["./neo-app"]