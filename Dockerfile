# ---- Build Stage ----
FROM --platform=linux/amd64 golang:1.24-alpine3.21 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Install go dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy only necessary source files
COPY *.go ./
COPY config/ ./config/
COPY controllers/ ./controllers/
COPY database/ ./database/
COPY middleware/ ./middleware/
COPY migrations/ ./migrations/
COPY models/ ./models/
COPY routes/ ./routes/
COPY utils/ ./utils/
COPY workers/ ./workers/

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main main.go

# ---- Final Stage ----
FROM alpine:3.21
WORKDIR /app

# Install minimal runtime dependencies
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the images directory (only if needed at runtime)
COPY images/ ./images/

# Copy the database content directory for blog seeding (only if needed at runtime)
COPY database/content/ ./database/content/

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8082

# Run the compiled binary
CMD ["./main"]