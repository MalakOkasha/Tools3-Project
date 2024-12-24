# Use a lightweight Go image to build the application
FROM golang:alpine AS build

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and source code into the container
COPY . .

# Install dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Use a minimal image to run the built Go app
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built application from the build stage
COPY --from=build /app/main .

# Ensure the binary has executable permissions
RUN chmod +x /app/main

# Set non-root user for OpenShift compatibility
RUN adduser -D -u 1001 appuser && chown -R appuser /app
USER appuser

# Expose the port your Go app runs on
EXPOSE 8080

# Start the Go app
CMD ["./main"]
