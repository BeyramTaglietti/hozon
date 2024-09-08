# Use the official Golang image from the Docker Hub
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o hozon .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# COPY env variables
COPY .env .env

# Define a build argument for the PostgreSQL version
ARG POSTGRES_VERSION=15

# Install pg_dump from the specified PostgreSQL client package
RUN apk --no-cache add "postgresql${POSTGRES_VERSION}-client"

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/hozon .

# Command to run the executable
ENTRYPOINT ["./hozon"]
