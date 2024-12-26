# Step 1: Build the binary using golang:alpine image
FROM golang:alpine AS build

# Set environment variables
ENV DB_USER="postgres01" \
    DB_PASSWORD="1234" \
    DB_HOST="db" \
    DB_PORT="5432" \
    DB_NAME="PTS"

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

# Ensure Go modules are installed
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Make sure the binary has execute permissions
RUN chmod +x main

# Step 2: Create the final image using alpine
FROM alpine:latest

# Set environment variables
ENV DB_USER="postgres01" \
    DB_PASSWORD="1234" \
    DB_HOST="db" \
    DB_PORT="5432" \
    DB_NAME="PTS"

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=build /app/main .

# Ensure the binary has execute permissions
RUN chmod +x ./main

# Expose the application port
EXPOSE 8080

# Set the command to run the binary
CMD ["./main"]

# Set labels for OpenShift build tracking
LABEL "io.openshift.build.commit.author"="Abdulrahman Hijazy <109446966+Abdulrahman-Hijazy@users.noreply.github.com>" \
      "io.openshift.build.commit.date"="Thu Dec 26 11:22:36 2024 +0200" \
      "io.openshift.build.commit.id"="5343a63314af2aff33555cc526fbc015247a639f" \
      "io.openshift.build.commit.message"="Update Dockerfile" \
      "io.openshift.build.commit.ref"="main" \
      "io.openshift.build.name"="back-9" \
      "io.openshift.build.namespace"="crt-20216056-dev" \
      "io.openshift.build.source-context-dir"="/Backend" \
      "io.openshift.build.source-location"="https://github.com/MalakOkasha/Tools3-Project.git"
