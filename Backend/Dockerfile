# Use a lightweight Go image
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

WORKDIR /root/
COPY --from=build /app/main .

# Expose the port your Go app runs on
EXPOSE 8080

# Start the Go app
CMD ["./main"]
