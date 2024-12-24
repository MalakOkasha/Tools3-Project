# Stage 1: Build the Go application
FROM golang:alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o main .

# Stage 2: Runtime image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary
COPY --from=build /app/main .

# Add a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /root && \
    chmod +x /root/main

# Switch to non-root user
USER 1004500000

EXPOSE 8080

CMD ["./main"]
