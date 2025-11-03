### I will run the projects always iside docker even for dev
## The very first time when the project is not created

1. Create a basic Dockerfile

```
# Base image with Go 1.24.5
FROM golang:1.24.5

# Set working directory
WORKDIR /app
```

Run

```
  docker run --rm -it -v $(pwd):/app neo:latest bash -c "go mod init github.com/carlosgab83/matrix/neo && \
  go get google.golang.org/grpc@latest && go mod tidy"
```

This creates an initial go.mod and go.sum

Then you have to use docker:

# Run in development
`docker run --rm -it -v $(pwd):/app -w /app neo:latest go run ./cmd/neo/main.go`

## Add a package (grpc@latest):
`docker run --rm -it -v $(pwd):/app neo:latest bash -c "go get google.golang.org/grpc@latest"`

## Generate gRPC code:
`docker run --rm -it -v $(pwd):/app neo:latest bash -c "cd internal/proto && protoc --go_out=. --go-grpc_out=. ingest_price.proto"`

### Build the image
`docker build -t neo:latest .`

### Run interactively, mounting local code for hot edits
`docker run --rm -it neo:latest`

### Run Go commands inside container
`docker run --rm -it -v $(pwd):/app neo:latest bash`
`docker run --rm -it -v $(pwd):/app neo:latest go run cmd/neo/main.go`

### Run Neo & Morpheus
`docker run --rm -it -v $(pwd):/app -w /app neo:latest bash cmd/nm.bash`




