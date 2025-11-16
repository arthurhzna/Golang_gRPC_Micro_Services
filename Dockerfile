# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the application source code
COPY . .

# Build the Go binary
RUN go build -o your_ip .

# Stage 2: Create a lightweight image with the Go binary
FROM alpine:3.18

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/your_ip your_ip

# Change ownership of the binary to the non-root user
RUN chown appuser:appgroup /app/your_ip

# Run the application as the non-root user
USER appuser

# Set the entry point command
CMD ["./your_ip"]