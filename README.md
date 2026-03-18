# Todos API

A production-ready RESTful API for managing todos with JWT authentication, built with Go, Gin, and PostgreSQL.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Usage Examples](#usage-examples)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [Security](#security)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Features

- ✅ **RESTful API** - Clean REST endpoints following best practices
- 🔐 **JWT Authentication** - Secure user authentication with JSON Web Tokens
- 📊 **PostgreSQL Database** - Reliable data persistence with GORM ORM
- 📝 **Swagger Documentation** - Interactive API documentation
- 🐳 **Docker Support** - Containerized application with Docker Compose
- ✨ **Graceful Shutdown** - Proper handling of shutdown signals
- 🧪 **Comprehensive Tests** - Unit and integration tests
- 🔧 **Environment Config** - Flexible configuration management
- 📦 **Database Migrations** - Automatic schema migrations with GORM
- 🚀 **Production Ready** - Optimized for deployment

## Architecture

```
todos/
├── config/          # Configuration management
├── database/        # Database connection and setup
├── handlers/        # HTTP request handlers
├── middleware/      # Custom middleware (auth, etc.)
├── models/          # Data models and schemas
├── routes/          # Route definitions
├── tests/           # Test files
├── docs/            # Swagger documentation (auto-generated)
├── .github/         # CI/CD workflows
├── main.go          # Application entry point
├── Dockerfile       # Docker configuration
├── docker-compose.yml
├── Makefile         # Build and development commands
└── README.md
```

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- Make (optional, for convenience commands)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/todos.git
cd todos
```

### 2. Setup environment variables

```bash
make env-setup
# Or manually: cp .env.example .env
# Edit .env with your configuration
```

Required environment variables:
```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=todos
DB_SSL_MODE=disable

JWT_SECRET=your_super_secret_jwt_key_change_in_production
JWT_EXPIRY_HOURS=24
```

### 3. Run with Docker (Recommended)

```bash
# Start all services (PostgreSQL + API)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### 4. Run locally

```bash
# Install dependencies
make install

# Start PostgreSQL (if not using Docker)
# Make sure PostgreSQL is running on localhost:5432

# Run the application
make dev
```

The API will be available at `http://localhost:8080`

## API Documentation

Once the server is running, visit:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

### API Endpoints

#### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register a new user | No |
| POST | `/api/v1/auth/login` | Login user | No |
| GET | `/api/v1/auth/profile` | Get user profile | Yes |

#### Todos

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/todos` | Get all todos | Yes |
| GET | `/api/v1/todos/:id` | Get a specific todo | Yes |
| POST | `/api/v1/todos` | Create a new todo | Yes |
| PUT | `/api/v1/todos/:id` | Update a todo | Yes |
| DELETE | `/api/v1/todos/:id` | Delete a todo | Yes |

## Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

### Create a todo (requires authentication)

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
```

### Get all todos

```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Development

### Available Make Commands

```bash
make help              # Show all available commands
make install           # Install dependencies
make dev              # Run in development mode
make build            # Build the binary
make run              # Build and run
make test             # Run tests
make test-coverage    # Run tests with coverage report
make clean            # Clean build artifacts
make docker-up        # Start Docker containers
make docker-down      # Stop Docker containers
make docker-rebuild   # Rebuild Docker containers
make swagger          # Generate Swagger docs
make lint             # Run linter
make format           # Format code
make watch            # Auto-reload on changes (requires air)
make env-setup        # Copy .env.example to .env
make db-reset         # Reset database (WARNING: deletes all data)
```

### Project Structure Details

#### Config Package
Centralized configuration management with validation and environment variable support.

#### Database Package
- Connection pooling (10 idle, 100 max connections)
- Automatic migrations with GORM
- Graceful connection closing

#### Middleware
- **Auth Middleware**: JWT token validation

#### Models
- User model with bcrypt password hashing
- Todo model with user relationship
- Input validation with Gin binding tags

#### Handlers
- Authentication (register, login, profile)
- Todo CRUD operations

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./tests/ -run TestRegister
```

### Test Configuration

Update `.env.test` with your test database credentials:
```env
APP_ENV=test
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=todos_test
JWT_SECRET=test_secret_key_for_testing_only
```

### Code Quality

```bash
# Format code
make format

# Run linter (requires golangci-lint)
make lint
```

## Deployment

### Docker Deployment

```bash
# Build and push to Docker Hub
docker build -t yourusername/todos-api:latest .
docker push yourusername/todos-api:latest

# Or use docker-compose
docker-compose up -d
```

### Environment Variables for Production

Make sure to set these in production:

```env
APP_ENV=production
JWT_SECRET=<strong-random-secret>
DB_SSL_MODE=require
```

Generate a secure JWT secret:
```bash
openssl rand -base64 32
```

### CI/CD

This project includes GitHub Actions workflow for:
- Linting with golangci-lint
- Running tests with coverage
- Building Docker images
- Code coverage reporting (Codecov)

Required GitHub Secrets:
- `DOCKER_USERNAME`
- `DOCKER_PASSWORD`

### Docker Image Optimizations

- Multi-stage build for smaller image size
- Non-root user execution for security
- Health check support
- Binary optimization with build flags

## Contributing

We welcome contributions! Here's how you can help:

### Development Setup

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make format`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Standards

- Follow Go best practices and [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Write tests for new features
- Update documentation as needed
- Keep functions small and focused
- Use meaningful variable names

### Commit Message Guidelines

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests when relevant

Example:
```
Add user authentication endpoint

Implement JWT-based authentication for user login and registration.
Closes #123
```

## Security

### Security Features

- ✅ Passwords hashed using bcrypt
- ✅ JWT tokens for stateless authentication
- ✅ SQL injection prevention with GORM
- ✅ Input validation on all endpoints
- ✅ Non-root Docker user
- ✅ Environment-based configuration

### Security Best Practices

- ⚠️ Add rate limiting for production
- ⚠️ Enable HTTPS in production
- ⚠️ Use strong JWT secrets (generate with `openssl rand -base64 32`)
- ⚠️ Keep dependencies updated
- ⚠️ Never commit secrets to version control
- ⚠️ Use environment variables for sensitive data

### Reporting Security Issues

If you discover a security vulnerability, please email security@yourproject.com instead of opening a public issue.

## Troubleshooting

### Database connection issues

```bash
# Check if PostgreSQL is running
docker-compose ps

# Check database logs
docker-compose logs postgres

# Reset database (WARNING: deletes all data)
make db-reset
```

### Port already in use

```bash
# Change APP_PORT in .env file
APP_PORT=3000
```

### Tests failing

```bash
# Make sure test database is configured
cp .env.example .env.test
# Edit .env.test with test database credentials

# Run tests
make test
```

### Swagger documentation not showing

```bash
# Regenerate Swagger docs
make swagger

# Restart the application
make dev
```

### Docker build fails

```bash
# Clean Docker cache and rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up
```

## Performance

- Database connection pooling configured (10 idle, 100 max)
- Connection lifetime management
- Optimized Docker image (multi-stage build)
- Binary size optimization with ldflags
- HTTP timeouts (10s read/write)

## Roadmap

- [ ] Add refresh tokens
- [ ] Implement pagination
- [ ] Add sorting and filtering
- [ ] Rate limiting middleware
- [ ] Email verification
- [ ] Password reset functionality
- [ ] WebSocket support for real-time updates
- [ ] Export todos to various formats
- [ ] Todo sharing and collaboration
- [ ] CORS middleware
- [ ] Request logging
- [ ] Metrics and monitoring

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Swag](https://github.com/swaggo/swag)

## Support

For issues and questions:
- Create an issue on GitHub
- Check existing documentation
- Review Swagger API docs at http://localhost:8080/swagger/index.html

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Version**: 1.0.0
**Last Updated**: March 2025
Made with ❤️ using Go
