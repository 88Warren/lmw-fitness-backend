#!/bin/sh

# Backend test runner for CI/CD pipeline
set -e

echo "🧪 Running backend unit tests..."

# Wait for database to be ready
echo "⏳ Waiting for database connection..."
for i in $(seq 1 30); do
    if go run -c 'package main; import ("database/sql"; _ "github.com/lib/pq"); func main() { db, err := sql.Open("postgres", "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable"); if err == nil { db.Ping() } }' 2>/dev/null; then
        echo "✅ Database connection established"
        break
    fi
    echo "⏳ Waiting for database... ($i/30)"
    sleep 2
done

# Download dependencies
echo "📦 Downloading Go modules..."
go mod download

# Run tests
echo "🚀 Executing tests..."
go test -v ./...

echo "✅ Backend tests completed successfully!"