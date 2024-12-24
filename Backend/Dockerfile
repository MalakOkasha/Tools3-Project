# Stage 1: Build the Go application
FROM golang:alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files (go.mod and go.sum) first
# This ensures Docker can cache the dependency installation step if go.mod or go.sum doesn't change.
COPY go.mod go.sum ./

# Download dependencies without modifying go.mod or go.sum
RUN go mod download

# Copy the rest of the application source code into the container
# This is done after installing dependencies to take advantage of Docker layer caching.
COPY . .

# Build the Go application
# The output binary will be named "main".
RUN go build -o main .

# Stage 2: Runtime image
# Use a minimal Alpine Linux image to run the built Go application
FROM alpine:latest

# Set the working directory for the runtime stage
WORKDIR /root/

# Copy the built Go application binary from the build stage to the runtime image
COPY --from=build /app/main .

# Ensure the binary has execution permissions
RUN chmod +x /root/main

# Expose the port your Go application listens on
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
