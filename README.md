# Chirpy: A Microblogging HTTP Server in Go

![Go](https://img.shields.io/badge/Go-00ADD8?style=plastic&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=plastic&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=plastic&logo=docker&logoColor=white)

A REST API backend for a microblogging platform, written in Go.

## What It Does

Chirpy provides a complete REST API for a microblogging application with:

- **User Management** - Create accounts, update profiles, secure password storage (Argon2)
- **Authentication** - JWT access tokens + refresh token rotation
- **Chirps** - Create, read, delete microblog posts (140 char max, profanity filter)
- **Webhooks** - Polka integration for premium user upgrades
- **Admin** - Metrics and database reset endpoints

## Quick Start

```bash
# Clone and install dependencies
git clone https://github.com/your-org/chirpy.git
cd chirpy
go mod download

# Start database
docker compose up -d
cp .env.example .env
# Edit .env with your secrets (JWT_SECRET, POLKA_KEY must be set)

# Run migrations
goose -s postgres "postgres://postgres:postgres@localhost/chirpy?sslmode=disable" up

# Run the server
go run .
```

## Configuration

Environment variables (required):

| Variable | Description |
|---------|-------------|
| `DB_URL` | PostgreSQL connection string |
| `PLATFORM` | Deployment environment (dev, prod) |
| `JWT_SECRET` | Secret key for JWT signing |
| `POLKA_KEY` | API key for Polka webhooks |

## Testing

```bash
go test -cover ./...
```

| Package       | Coverage |
|---------------|----------|
| main          | 1.7%     |
| internal/auth | 70.2%    |

Tests use Go's table-driven test pattern.

## API Endpoints

### Health & Admin

```bash
# Health check
GET /api/healthz
# Response: {"status":"ok"}

# Metrics (fileserver hits)
GET /admin/metrics
# Response: {"hits": 42}

# Reset database
POST /admin/reset
# Response: {"success": true}
```

### Users

```bash
# Create user
POST /api/users
Content-Type: application/json
{"email": "alice@example.com", "password": "securepassword123"}
# Response: {"id": "uuid", "email": "alice@example.com", "is_chirpy_red": false}

# Update user (auth required)
PUT /api/users
Authorization: Bearer <jwt>
Content-Type: application/json
{"email": "newemail@example.com", "password": "newpassword"}
# Response: {"id": "uuid", "email": "newemail@example.com", "is_chirpy_red": false}
```

### Authentication

```bash
# Login
POST /api/login
Content-Type: application/json
{"email": "alice@example.com", "password": "securepassword123"}
# Response: {"id": "...", "email": "alice@example.com", "is_chirpy_red": false, "token": "...", "refresh_token": "..."}

# Refresh token
POST /api/refresh
Authorization: Bearer <refresh_token>
# Response: {"token": "new_access_token"}

# Revoke token
POST /api/revoke
Authorization: Bearer <refresh_token>
# Response: {"success": true}
```

### Chirps

```bash
# Get all chirps
GET /api/chirps?author_id=<uuid>&sort=asc|desc
# Query params:
#   - author_id: filter by user (optional)
#   - sort: sort by creation time, "asc" or "desc" (default: asc)
# Response: [{"id": "...", "body": "Hello world", "user_id": "...", ...}]

# Get chirp by ID
GET /api/chirps/{chirpID}
# Response: {"id": "...", "body": "Hello world", "user_id": "...", ...}

# Create chirp (auth required)
POST /api/chirps
Authorization: Bearer <jwt>
Content-Type: application/json
{"body": "My first chirp!"}
# Response: {"id": "...", "body": "My first chirp!", "user_id": "...", ...}

# Delete chirp (auth required, owner only)
DELETE /api/chirps/{chirpID}
Authorization: Bearer <jwt>
# Response: {"success": true}
```

### Webhooks

```bash
# Polka webhooks
POST /api/polka/webhooks
Authorization: Polka <polka_key>
Content-Type: application/json
{"event": "user.upgraded", "data": {"user_id": "..."}}
# Response: {"success": true}
```

## Development

```bash
# Available commands (via Justfile)
just run          # Run server
just build        # Build binary
just test         # Run tests
just lint         # Lint code
just tidy         # Update dependencies
just db           # Start database
just bruno       # Run API tests (Bruno)
```

## Project Structure

```
chirpy/
├── main.go           # Entry point, router setup
├── handler_*.go       # HTTP handlers
├── internal/
│   ├── auth/         # JWT, password hashing
│   └── database/     # SQLC generated
├── sql/
│   ├── schema/      # Database migrations
│   └── queries/     # Custom SQL
├── bruno/           # API tests
└── justfile         # Development tasks
```

## Tech Stack

- **Language**: Go 1.25+
- **Database**: PostgreSQL 16+
- **ORM**: SQLC (compile-time type-safe SQL)
- **Auth**: JWT (golang-jwt), Argon2 password hashing
- **Migrations**: goose