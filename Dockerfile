# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Install go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application for Ionis VPS
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

# ---- Final Stage ----
FROM alpine:latest
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the images directory
COPY --from=builder /app/images ./images

# Expose the application port
EXPOSE 8082

# Run the compiled binary
CMD ["./main"]