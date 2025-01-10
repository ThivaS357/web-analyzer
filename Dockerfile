# Start with a lightweight Go image
FROM golang:1.22-alpine AS build

# Set environment variables for Go
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go binary
RUN go build -o web-analyzer main.go

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder
COPY --from=build /app/web-analyzer .

# Copy HTML templates
COPY templates ./templates

# Expose the port the app runs on
EXPOSE 8080

# Command to run the binary
CMD ["./web-analyzer"]
