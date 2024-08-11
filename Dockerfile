FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /build

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Set necessary environment variables for building the Go application
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go application
RUN go build -ldflags="-s -w" -o apiserver cmd/api/main.go

# Stage 2: Create a minimal image with the built binary
FROM alpine:latest

# Copy the binary and necessary files from the builder stage
COPY --from=builder /build/apiserver /apiserver
COPY --from=builder /build/.env /
COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set the timezone
ENV TZ Asia/Jakarta

# Command to run when starting the container
ENTRYPOINT ["/apiserver"]