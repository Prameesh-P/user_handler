# Use an official Golang runtime as the base image
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules file and download dependencies
COPY go.mod go.sum ./
# RUN go mod download
RUN go mod download
# Run go mod tidy to update and clean dependencies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application

RUN go build -o app cmd/main.go

# Use a minimal base image for the final runtime
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose the port for the Golang application
EXPOSE 8081

# Start the Golang application
CMD ["./app"]