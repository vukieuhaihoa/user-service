# User Service

A microservice for user management built with **Go**, following **Domain-Driven Design (DDD)** and **Clean Architecture** principles. Part of the Bookmark Management DDD system.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25+ |
| Web Framework | Gin |
| Database | PostgreSQL + GORM |
| Cache | Redis |
| Authentication | JWT (RSA-2048) |
| API Docs | Swagger / OpenAPI |
| Logging | Zerolog |
| Containerization | Docker (multi-stage) |
| CI/CD | GitHub Actions |

---

## Project Structure

```
user-service/
├── cmd/
│   ├── api/main.go          # API server entry point
│   └── migrate/main.go      # Database migration entry point
├── internal/
│   ├── api/                 # Gin engine setup, routing, middleware
│   ├── app/
│   │   ├── handler/         # HTTP request handlers
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Data access layer
│   │   └── model/           # Domain models
│   ├── infrastructure/      # Dependency injection, DB/Redis/JWT init
│   └── docs/                # Generated Swagger documentation
├── migrations/              # SQL migration files
├── postgres/init_db/        # DB initialization scripts
├── .github/workflows/       # CI/CD pipelines
├── Dockerfile
├── docker-compose.dev.yaml
└── Makefile
```

---

## API Endpoints

### Public

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health-check` | Service health check |
| `POST` | `/v1/users/register` | Register a new user |
| `POST` | `/v1/users/login` | Login and receive JWT |
| `GET` | `/swagger/*` | Swagger UI |

### Protected (JWT required)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/v1/self/info` | Get current user profile |
| `PUT` | `/v1/self/info` | Update current user profile |

> Include the JWT token in the `Authorization: Bearer <token>` header for protected routes.

---

## Getting Started

### Prerequisites

- [Go 1.25+](https://golang.org/)
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/)

### 1. Start infrastructure (PostgreSQL + Redis)

```bash
make dev-up
```

### 2. Run database migrations

```bash
make migrate
```

### 3. Generate Swagger docs and run the server

```bash
make dev-run
```

The API will be available at `http://localhost:8080`.
Swagger UI: `http://localhost:8080/swagger/index.html`

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `:8080` | HTTP server port |
| `SERVICE_NAME` | `user-service` | Service name for logging |
| `APP_HOST_NAME` | `localhost:8080` | Host used in Swagger docs |
| `INSTANCE_ID` | *(random UUID)* | Unique instance identifier |

---

## Development

### Run tests

```bash
make test
```

> Requires 80% code coverage to pass.

### Generate mocks

```bash
make mock-gen
```

### Generate Swagger docs

```bash
make swag-gen
```

### Generate RSA keys for JWT

```bash
make generate-rsa-key
```

### Create a new migration

```bash
make new-schema name=<migration_name>
```

---

## Docker

### Build image

```bash
make docker-build
```

### Run tests in Docker

```bash
make docker-test
```

### Push image to Docker Hub

```bash
make docker-release
```

---

## CI/CD

### CI Pipeline

Triggered on:
- Push to `main`
- Pull requests targeting `main`
- Version tags (`v*.*.*`)

Steps: run tests → build Docker image → push to Docker Hub

### CD Pipeline

Triggered when CI passes on a non-`main` branch. Deploys to a self-hosted runner by pulling the latest image and restarting the service.

---

## Database

### Schema

```sql
CREATE TABLE users (
  id           varchar(36) PRIMARY KEY,
  display_name varchar(255)  NOT NULL,
  username     varchar(255)  NOT NULL UNIQUE,
  password     varchar(2048) NOT NULL,
  email        varchar(2048) NOT NULL UNIQUE,
  created_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at   TIMESTAMPTZ   -- soft delete
);
```

### Run migrations manually

```bash
make migrate
```
