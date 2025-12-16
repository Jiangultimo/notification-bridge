# Notification Bridge

Notification Bridge API Server

## Prerequisites

- Go 1.25.5+
- [Air](https://github.com/air-verse/air) - Hot reload for Go apps
- [Swag](https://github.com/swaggo/swag) - Swagger documentation generator

### Install Development Tools

```bash
# Install Air (hot reload)
go install github.com/air-verse/air@latest

# Install Swag CLI (swagger docs generator)
go install github.com/swaggo/swag/cmd/swag@latest
```

## Getting Started

### Install Dependencies

```bash
go mod download
```

### Run Development Server (with hot reload)

```bash
air
```

The server will start at `http://localhost:8000`

### Run Without Hot Reload

```bash
go run ./cmd/server
```

## API Documentation

Swagger UI is available at: `http://localhost:8000/swagger/`

### Regenerate Swagger Docs

When you modify API annotations, regenerate the docs:

```bash
swag init -g ./cmd/server/main.go -o ./docs
```
