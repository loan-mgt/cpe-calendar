# Build stage
FROM golang:1.23-alpine AS builder

# Create a working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o /app/calendar-app main.go

# Final stage
FROM alpine:latest

# Create a working directory
WORKDIR /root/

# Copy the compiled Go binary from the build stage
COPY --from=builder /app/calendar-app .

# Copy static files to the container
COPY static ./static
COPY secret ./secret

# Expose the port on which the Go app will run
EXPOSE 8080

# Command to run the Go app
CMD ["./calendar-app"]
