# GoFiber Boilerplate

A professional-grade Go Fiber boilerplate with clean architecture, repository patterns, PASETO authentication, RBAC, MongoDB sharded transactions, Redis, and robust unit testing.

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
│  │                 │  │    (with Atomic Transactions)   │   │
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
- **Repository Pattern**: Abstracted data access for testability (includes full Mock support)
- **Dependency Injection**: All dependencies injected through constructors
- **SOLID Principles**: Focused on maintainability and interface segregation
- **Atomic Operations**: Support for MongoDB Transactions in the Service layer

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
│   │   │   └── mock/        # Repository mocks for unit testing [NEW]
│   │   ├── service/         # User business logic
│   │   └── routes.go        # User routes registration
│   ├── database/            # Database operations
│   │   └── migration/       # Centralized index migrations [NEW]
│   ├── middleware/          # Cross-cutting middleware
│   │   ├── auth.go          # PASETO authentication
│   │   └── rbac.go          # Role-based access control
│   └── config/              # Configuration management
├── shared/
│   └── entity/              # Shared domain entities
├── pkg/
│   ├── database/            # MongoDB & Redis connections
│   ├── logger/              # Zap logger setup
│   ├── response/            # API response helpers
│   ├── token/               # PASETO token maker
│   ├── utils/               # Performance-optimized helpers [NEW]
│   └── validator/           # Request validation
├── Makefile                 # Development commands [NEW]
├── air.toml                 # Live-reloading config [NEW]
├── docker-compose.yml       # MongoDB sharded + Redis
└── Dockerfile               # Multi-stage build
```

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [Air](https://github.com/air-verse/air) (optional, for live reloading)

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

4. **Run with Live Reload (Recommended)**
   ```bash
   make dev
   # or simply 'air'
   ```

5. **Run standard execution**
   ```bash
   make run
   # or 'go run cmd/api/main.go'
   ```

## 🔄 Core Features

### 📦 Database Migrations
Indexing is centralized in `internal/database/migration`. On every startup, the application verifies and creates necessary MongoDB indexes, ensuring data integrity without manual intervention.

### ⚛️ Atomic Transactions
The Service layer supports MongoDB Transactions. The `RegisterWithStats` method in `UserService` demonstrates atomic cross-collection updates:
- Ensures user creation and stats initialization happen together.
- Automatic rollback on any failure within the transaction block.

### 🧪 Unit Testing with Mocks
Isolation is key. We provide a manual mocking system for repositories:
- **Location**: `internal/user/repository/mock/`
- **Example Usage**: See `internal/user/service/user_service_test.go`
```bash
go test -v ./...
```

### 🇮🇩 Indonesian Time Helper
Localized time formatting in `pkg/utils/time.go`:
- Formats: `Sabtu, 27 Desember 2025 - 22:15 WIB`
- Optimized: Timezone location (`Asia/Jakarta`) is loaded once at package level.

## 🔐 Authentication

### PASETO Tokens
Uses **PASETO V2 (local)** which is more secure than standard JWT as it avoids algorithm confusion attacks.
- **Access Token**: Short-lived (15 min default)
- **Refresh Token**: Stored in Redis for session management

### Endpoints
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | Login | ❌ |
| POST | `/api/v1/auth/refresh` | Refresh tokens | ❌ |
| POST | `/api/v1/auth/logout` | Logout | ✅ |

## 👤 User Management

### RBAC Roles
- **ADMIN**: Full access
- **USER**: Access to own profile

### Filtering & Searching (ADMIN)
The `GET /api/v1/users` endpoint supports advanced dynamic filtering using **kebab-case** parameters:

| Parameter | Field Map | Description |
|-----------|-----------|-------------|
| `first-name` | `firstName` | Filter by first name |
| `last-name` | `lastName` | Filter by last name |
| `email` | `email` | Filter by email |
| `role` | `role` | Filter by role (`ADMIN`/`USER`) |
| `is-active` | `isActive` | Filter by status (`true`/`false`) |
| `search` | N/A | Global regex search across name and email |

### Endpoints
| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| GET | `/api/v1/users/me` | Get current user | USER |
| GET | `/api/v1/users` | List all users (with filters) | ADMIN |
| GET | `/api/v1/users/:id` | Get user by ID | ADMIN |
| DELETE | `/api/v1/users/:id` | Delete user | ADMIN |

## 📦 MongoDB Sharded Cluster
The `docker-compose.yml` sets up a complete sharded cluster with a Query Router (**mongos**), demonstrating production-ready horizontal scaling patterns.

## 🔄 CI/CD Pipeline
- **Lint**: golangci-lint
- **Test**: Unit tests with coverage
- **Security**: govulncheck vulnerability scanning

## 📝 API Response Format

### Success Response
```json
{
  "success": true,
  "code": 201,
  "status": "CREATED",
  "message": "user registered successfully",
  "data": {
    "createdAt": "Sabtu, 27 Desember 2025 - 22:15 WIB",
    "updatedAt": "Sabtu, 27 Desember 2025 - 22:15 WIB"
  }
}
```

### Paginated Response (Laravel-style)
```json
{
  "success": true,
  "code": 200,
  "status": "OK",
  "data": [...],
  "meta": {
    "currentPage": 1,
    "from": 1,
    "lastPage": 5,
    "links": [
      { "url": null, "label": "&laquo; Previous", "active": false },
      { "url": "http://localhost:3000/api/v1/users?per-page=10&page=1", "label": "1", "active": true },
      { "url": "http://localhost:3000/api/v1/users?per-page=10&page=2", "label": "Next &raquo;", "active": false }
    ],
    "path": "http://localhost:3000/api/v1/users",
    "perPage": 10,
    "to": 10,
    "total": 50,
    "firstPageUrl": "http://localhost:3000/api/v1/users?page=1",
    "lastPageUrl": "http://localhost:3000/api/v1/users?page=5",
    "nextPageUrl": "http://localhost:3000/api/v1/users?page=2",
    "prevPageUrl": null
  }
}
```

## 📜 License
MIT License
