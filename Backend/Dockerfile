# Stage 1: Build the Go application
FROM golang:alpine AS build

# Set the working directory in the container
WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built application binary
COPY --from=build /app/main .

# Ensure the binary has the proper execution permissions
RUN chmod +x /root/main

# Add a non-root user and group for OpenShift
RUN addgroup -g 1004500000 appgroup && \
    adduser -S -u 1004500000 -G appgroup appuser && \
    chown -R appuser:appgroup /root

# Switch to a non-root user (required for OpenShift)
USER 1004500000

# Expose the port
EXPOSE 8080

# Start the Go app
CMD ["./main"]
