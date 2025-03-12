# ---- Build Stage ----
    FROM golang:1.23-alpine AS builder
    WORKDIR /app
    
    # Install dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy source code
    COPY . .
    
    # Build the Go application
    RUN go build -o main main.go
    
    # ---- Final Stage ----
    FROM alpine:latest
    WORKDIR /app
    
    # Copy only the built binary from the builder stage
    COPY --from=builder /app/main .
    
    # Expose the application port
    EXPOSE 8082
    
    # Run the compiled binary
    CMD ["./main"]