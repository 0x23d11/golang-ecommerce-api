# Use an official Go runtime as a parent image
# Choose a version that matches your local Go version (e.g., 1.22)
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies.
# This leverages Docker layer caching. Dependencies are only downloaded again
# if go.mod or go.sum change.
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 produces a statically linked binary suitable for minimal base images
# -o /app/server builds the output binary named 'server' inside the /app directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/api/main.go


# --- Second Stage: Create the final, smaller image ---

# Use a minimal base image like Alpine Linux
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server /app/server

# Copy the .env file if you are using one for configuration
# COPY .env .env

# Expose the port the app runs on (matches SERVER_PORT in config/env)
# Make sure this matches the port you intend to use (default is 8080)
EXPOSE 8080

# Command to run the executable
# The binary name 'server' must match the one specified in the build step
CMD ["/app/server"]