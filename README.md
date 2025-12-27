# GoFiber Boilerplate

A production-ready Go Fiber boilerplate with clean architecture, repository pattern, PASETO authentication, RBAC, MongoDB sharded cluster, Redis, and CI/CD.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      API Layer (Fiber)                       │
│  ┌─────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │ Request │→ │  Middleware │→ │        Handlers         │  │
│  └─────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│  ┌─────────────────┐  ┌─────────────────────────────────┐   │
│  │   AuthService   │  │          UserService            │   │
│  └─────────────────┘  └─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Repository Layer                          │
│  ┌─────────────────┐  ┌─────────────────────────────────┐   │
│  │  UserRepository │  │       TokenRepository           │   │
│  │    (MongoDB)    │  │          (Redis)                │   │
│  └─────────────────┘  └─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Design Principles

- **Clean Architecture**: Separation of concerns with distinct layers
- **Repository Pattern**: Abstracted data access for testability
- **Dependency Injection**: All dependencies injected through constructors
- **SOLID Principles**: Single responsibility, interface segregation

## 📁 Project Structure (Feature-First)

```
gofiber-boilerplate/
├── cmd/api/
│   └── main.go              # Application entry point
├── internal/
│   ├── auth/                # Auth feature module
│   │   ├── dto/             # Auth DTOs
│   │   ├── handler/         # Handlers (login, register, refresh, logout)
│   │   ├── repository/      # Token repository (Redis)
│   │   ├── service/         # Auth business logic
│   │   └── routes.go        # Auth routes registration
│   ├── user/                # User feature module
│   │   ├── dto/             # User DTOs
│   │   ├── handler/         # Handlers (me, list, update)
│   │   ├── repository/      # User repository (MongoDB)
│   │   ├── service/         # User business logic
│   │   └── routes.go        # User routes registration
│   ├── middleware/          # Cross-cutting middleware
│   │   ├── auth.go          # PASETO authentication
│   │   └── rbac.go          # Role-based access control
│   └── config/              # Configuration management
├── shared/
│   └── entity/              # Shared domain entities
│       └── user.go          # User entity (shared across features)
├── pkg/
│   ├── database/            # MongoDB & Redis connections
│   ├── logger/              # Zap logger setup
│   ├── response/            # API response helpers
│   ├── token/               # PASETO token maker
│   └── validator/           # Request validation
├── .github/workflows/       # CI/CD pipeline
├── docker-compose.yml       # MongoDB sharded + Redis
├── Dockerfile               # Multi-stage build
└── README.md
```

### Feature-First Design Benefits

- **Cohesion**: All auth code in `internal/auth/`, all user code in `internal/user/`
- **Self-contained features**: Each feature has its own routes.go for registration
- **Easy to scale**: Add new features as separate modules
- **Clear boundaries**: Features can evolve independently

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/itsahyarr/gofiber-boilerplate.git
   cd gofiber-boilerplate
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Start dependencies**
   ```bash
   docker-compose up -d redis mongo-config mongo-shard1 mongo-shard2 mongos mongo-init
   ```

4. **Wait for MongoDB initialization** (first time only)
   ```bash
   docker-compose logs -f mongo-init
   # Wait until you see "MongoDB Sharded Cluster initialized successfully!"
   ```

5. **Install Go dependencies**
   ```bash
   go mod download
   ```

6. **Run the application**
   ```bash
   go run ./cmd/api
   ```

### Docker Deployment

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f api
```

## 🔐 Authentication

### PASETO Tokens

This boilerplate uses **PASETO V2 (local)** for secure, stateless authentication:

- **Access Token**: Short-lived (15 min default), used for API requests
- **Refresh Token**: Long-lived (7 days), stored in Redis for management

### Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | Login | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh tokens | ❌ |
| POST | `/api/v1/auth/logout` | Logout | ✅ |

### Example: Register

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Example: Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Using Access Token

```bash
curl http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer <access_token>"
```

## 👤 User Management

### RBAC Roles

- **ADMIN**: Full access to all resources
- **USER**: Access to own profile only

### Endpoints

| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| GET | `/api/v1/users/me` | Get current user | USER |
| PUT | `/api/v1/users/me/password` | Change password | USER |
| GET | `/api/v1/users` | List all users | ADMIN |
| GET | `/api/v1/users/:id` | Get user by ID | ADMIN |
| PUT | `/api/v1/users/:id` | Update user | ADMIN/Self |
| DELETE | `/api/v1/users/:id` | Delete user | ADMIN |

## 🔧 Configuration

All configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | development | Environment (development/production) |
| `LOG_LEVEL` | debug | Log level (debug/info/warn/error) |
| `SERVER_PORT` | 3000 | Server port |
| `MONGODB_URI` | mongodb://localhost:27017 | MongoDB connection URI |
| `MONGODB_DATABASE` | gofiber_boilerplate | MongoDB database name |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `TOKEN_SYMMETRIC_KEY` | - | 32-char PASETO key |
| `ACCESS_TOKEN_DURATION` | 15m | Access token lifetime |
| `REFRESH_TOKEN_DURATION` | 168h | Refresh token lifetime |

## 📦 MongoDB Sharded Cluster

The `docker-compose.yml` sets up a complete MongoDB sharded cluster:

- **Config Server**: 1 replica set member
- **Shard 1**: 1 replica set member
- **Shard 2**: 1 replica set member
- **Mongos Router**: Query router

This setup is optimized for single-server deployments while demonstrating sharding patterns.

## 🔄 CI/CD Pipeline

The GitHub Actions workflow includes:

1. **Lint**: golangci-lint
2. **Test**: Unit tests with coverage
3. **Security**: govulncheck vulnerability scanning
4. **Build**: Binary and Docker image
5. **Deploy**: Staging (develop) / Production (main)

## 📝 API Response Format

### Success Response

```json
{
  "success": true,
  "message": "user retrieved successfully",
  "data": {
    "id": "...",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "USER"
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "invalid email or password",
  "error": {
    "code": "UNAUTHORIZED",
    "details": ""
  }
}
```

### Paginated Response

```json
{
  "success": true,
  "message": "users retrieved successfully",
  "data": [...],
  "pagination": {
    "current_page": 1,
    "page_size": 10,
    "total_items": 100,
    "total_pages": 10
  }
}
```

## 🧪 Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -cover ./...

# Run specific package tests
go test -v ./internal/service/...
```

## 📊 Observability

### Zap Logger

Structured logging with:
- Environment-aware formatting (JSON for production, colorized for dev)
- Configurable log levels
- Request ID tracking
- Error context capture

### Health Check

```bash
curl http://localhost:3000/health
```

## 🔒 Security Features

- PASETO tokens (more secure than JWT)
- Password hashing with bcrypt
- Role-based access control (RBAC)
- Redis-backed refresh token management
- Non-root Docker container
- Graceful shutdown

## 📜 License

MIT License
