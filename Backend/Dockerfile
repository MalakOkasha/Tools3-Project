# Use a lightweight Go image for the build phase
FROM golang:alpine AS build

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and source code into the container
COPY . .

# Install dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Ensure the binary has the correct permissions during the build phase
RUN chmod +x main
RUN ls -l main  # Check permissions of the built binary

# Use a minimal image to run the built Go app
FROM alpine:latest

# Set the working directory for the runtime container
WORKDIR /root/

# Copy the built application from the build container
COPY --from=build /app/main .

# Ensure the binary has the correct permissions in the runtime container
RUN chmod +x ./main
RUN ls -l ./main  # Check permissions of the copied binary

# Expose the port your Go app runs on
EXPOSE 8080

# Start the Go application
CMD ["./main"]
