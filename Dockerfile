# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Set environment variables to disable CGO and ensure static linking
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go binary
RUN go build -o hello .

# Stage 2: Create the runtime image using Alpine
FROM alpine:latest

# Install required dependencies (e.g., certificates)
RUN apk --no-cache add ca-certificates

# Set the working directory to /root
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/hello /root/hello

# Expose port 8080
EXPOSE 8080

# Run the Go binary
CMD ["/root/hello"]
