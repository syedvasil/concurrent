# Build stage
FROM golang:1.22.1 AS builder

# Set the working directory inside the container
WORKDIR /app

RUN pwd && ls

# Copy the Go modules manifests
#COPY go.mod go.sum ./
COPY go.mod ./
# Download the Go modules dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

RUN echo "Contents after copying go.mod and go.sum:" && pwd && ls -al /app


# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install necessary CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/cities.json /root/
COPY --from=builder /app/tmp /root/tmp

# Ensure the binary has execute permissions (optional if already executable)
RUN chmod +x ./main

RUN echo "Contents after copying the Go binary:" && pwd && ls -al

# Expose the port that the application runs on
EXPOSE 8080

# Command to run the Go binary
CMD ["./main"]
