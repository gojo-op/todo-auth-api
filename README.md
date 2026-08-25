# Todo Auth API

REST API for user authentication and todo management built with Go, Gin, GORM, and SQLite.

## Features

- User registration and login with JWT
- Password hashing with bcrypt
- User-scoped todo CRUD
- Environment-based configuration
- Graceful server shutdown

## Project Structure

```
cmd/api/              Application entrypoint
internal/config/      Environment configuration
internal/database/    Database initialization and migrations
internal/handler/     HTTP handlers
internal/middleware/  JWT authentication middleware
internal/models/      Data models
internal/repository/  Database access layer
internal/routes/      Route registration
internal/service/     Business logic
```

## Setup

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Install dependencies:

```bash
go mod download
```

3. Run the API:

```bash
go run ./cmd/api
```

For hot reload during development:

```bash
air
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SQLITE_PATH` | Yes | - | SQLite database file path |
| `SERVER_PORT` | No | `8080` | HTTP server port |
| `JWT_SECRET` | Yes | - | Secret used to sign JWT tokens |

## API Endpoints

### Health

- `GET /ping`

### Auth

- `POST /api/auth/register`
- `POST /api/auth/login`

### Todos

All todo routes require `Authorization: Bearer <token>`.

- `GET /api/todos`
- `POST /api/todos`
- `GET /api/todos/:id`
- `PUT /api/todos/:id`
- `DELETE /api/todos/:id`

## Example Requests

Register:

```bash
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane","email":"jane@example.com","password":"secret123"}'
```

Create todo:

```bash
curl -X POST http://localhost:8000/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Buy milk","description":"2 liters"}'
```

## Tests

```bash
go test ./...
```

## Notes

- Do not commit `.env`, database files, or build artifacts.
- Use `.env.example` as the template for local configuration.
- Change `JWT_SECRET` before deploying to production.
