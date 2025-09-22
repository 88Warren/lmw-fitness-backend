# LMW Fitness Backend API

A Go-based REST API for the LMW Fitness platform, providing workout programs, payment processing, and user management.

## Features

- **Payment Processing**: Stripe integration for workout program purchases
- **Email Automation**: Brevo integration for newsletters and transactional emails
- **Workout Management**: CRUD operations for fitness programs and workouts
- **User Management**: User registration and profile management
- **Background Jobs**: Asynchronous payment processing
- **Health Monitoring**: Comprehensive health checks and metrics

## Tech Stack

- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL with GORM
- **Logging**: Zap structured logging
- **Containerization**: Docker
- **Deployment**: Kubernetes with ArgoCD

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 13+
- Docker (optional)

### Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository>
   cd backend
   cp .env.development.example .env.development
   ```

2. **Configure environment**:
   Edit `.env.development` with your local settings

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Run database migrations**:
   ```bash
   go run main.go
   ```

5. **Start the server**:
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8082`

### Testing

```bash
# Run all tests
go test ./tests/...

# Run tests with coverage
go test -cover ./tests/...

# Run specific test
go test -run TestHealthEndpoint ./tests/
```

### Docker Development

```bash
# Build and run with Docker Compose
docker-compose up --build

# Run in background
docker-compose up -d
```

## API Documentation

### Health Endpoints

- `GET /monitoring/health` - Comprehensive health check with metrics
- `GET /monitoring/ready` - Kubernetes readiness probe
- `GET /monitoring/live` - Kubernetes liveness probe

### Core Endpoints

- `POST /api/create-checkout-session` - Create Stripe checkout session
- `POST /api/stripe-webhook` - Handle Stripe webhooks
- `GET /api/workouts` - Get workout programs
- `POST /api/newsletter/subscribe` - Newsletter subscription

### Authentication

Currently using session-based authentication. JWT implementation planned.

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GO_ENV` | Environment (development/production/test) | Yes |
| `PORT` | Server port | No (default: 8082) |
| `DB_HOST` | Database host | Yes |
| `DB_USER` | Database user | Yes |
| `DB_PASSWORD` | Database password | Yes |
| `DB_NAME` | Database name | Yes |
| `STRIPE_SECRET_KEY` | Stripe secret key | Yes |
| `BREVO_API_KEY` | Brevo API key | Yes |

## Deployment

### Kubernetes

The application is deployed using ArgoCD with GitOps workflow:

1. **Build**: Docker images are built and pushed to registry
2. **Deploy**: ArgoCD syncs from Git repository
3. **Monitor**: Health checks ensure service availability

### Health Checks

- **Liveness**: `/monitoring/live` - Basic application health
- **Readiness**: `/monitoring/ready` - Database connectivity check
- **Detailed**: `/monitoring/health` - Full system metrics

## Monitoring

### Structured Logging

All requests are logged with structured data:
- Request method, path, status
- Response time and body size
- Client IP and user agent
- Error details for failed requests

### Metrics Collection

- Request counts by endpoint
- Response time tracking
- Error rate monitoring
- System resource usage

### Alerts

Configure alerts for:
- High error rates (>5%)
- Slow responses (>2s average)
- Database connection failures
- Memory usage >80%

## Contributing

1. Create feature branch from `main`
2. Write tests for new functionality
3. Ensure all tests pass
4. Update documentation
5. Submit pull request

## Security

- Environment variables for sensitive data
- CORS middleware configured
- Input validation on all endpoints
- Rate limiting (planned)
- SQL injection protection via GORM

## Performance

- Database connection pooling
- Background job processing
- Structured logging for debugging
- Health check caching (planned)
