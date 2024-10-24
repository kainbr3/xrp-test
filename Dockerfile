# Use the official Golang image to create a build artifact.
# This is the build stage.
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates openssl

# Set Certificate Location
ARG cert_location=/usr/local/share/ca-certificates

# Get certificate from "github.com"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt

# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/proxy.golang.crt

# Get certificate from "storage.googleapis.com"
RUN openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/storage.googleapis.crt

# Update certificates
RUN update-ca-certificates

# Set Working Directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy security files
COPY global-bundle.pem secret.key ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# List all files on source dir
RUN ls -la

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Start a new stage from scratch
FROM golang:1.22

# Label and Install Bash
LABEL maintainer="getBRAZA"
RUN apt-get update && apt-get install -y bash

# Set Working Directory
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the security files from the previous stage
COPY --from=builder /app/global-bundle.pem /app/global-bundle.pem
COPY --from=builder /app/secret.key /app/secret.key

# List all files on destination dir
RUN ls -la

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]