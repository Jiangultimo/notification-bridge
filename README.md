# Notification Bridge

A lightweight webhook notification bridge server built with Go. Receives webhooks from NetAlertX (for now) and forwards formatted notifications to services like Bark.

## Features

- 🔔 **Webhook Reception**: Accept webhooks from NetAlertX Webhook
- 📤 **Notification Forwarding**: Forward to Bark (iOS push notifications)
- 🎨 **Markdown Formatting**: Automatically format webhook data into readable Markdown
- ⚡ **Async Processing**: Non-blocking notification delivery
- 📝 **Request Tracking**: Built-in request ID for debugging
- 📊 **API Documentation**: Interactive Swagger UI

## Supported Integrations

- **Sources**: NetAlertX (network monitoring)
- **Destinations**: Bark (iOS push notifications)

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

## Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Bark Configuration (iOS Push Notifications)
BARK_API_URL=https://your-bark-server.com
BARK_KEY=your-bark-key
BARK_ICON=https://your-icon-url.png  # Optional

# Server Configuration
PORT=8000  # Optional, defaults to 8000
```

**Required:**
- `BARK_API_URL`: Your Bark server URL
- `BARK_KEY`: Your Bark device key

**Optional:**
- `BARK_ICON`: Custom icon URL for notifications
- `PORT`: Server port (defaults to 8000)

## Getting Started

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment

Create a `.env` file with your Bark credentials (see [Configuration](#configuration))

### 3. Run Development Server (with hot reload)

```bash
air
```

The server will start at `http://localhost:8000`

### 4. Run Without Hot Reload

```bash
go run ./cmd/server
```

## Usage

### Available Endpoints

- **Health Check**: `GET /healthz`
  - Returns server health status

- **NetAlertX Webhook**: `POST /api/webhook/netalertx`
  - Receives NetAlertX notifications and forwards to Bark

- **Swagger UI**: `GET /swagger/`
  - Interactive API documentation

### Example: Configure NetAlertX

In your NetAlertX settings, configure the webhook URL:

```
http://your-server:8000/api/webhook/netalertx
```

NetAlertX will send notifications about:
- 🆕 New devices discovered on your network
- 🔁 Devices reconnecting after being offline
- 🔴 Devices going offline
- 🔌 Plugin events

These will be automatically formatted and sent to your Bark app.

### Notification Format

Webhook data is automatically formatted into Markdown for Bark notifications. Here's an example of what you'll receive:

```markdown
**Notification Title:** `NetAlertX`

**Message Content:**
### 🆕 新设备 (1)
- **设备**: Xiaomi Phone
  - **MAC**: aa:bb:cc:dd:ee:ff
  - **IP**: 192.168.1.100
  - **时间**: 2025-01-19 10:30:00+08:00

### 🔁 重连设备 (1)
- **设备**: MacBook Pro
  - **MAC**: 11:22:33:44:55:66
  - **IP**: 192.168.1.50
  - **厂商**: Apple Inc.
  - **时间**: 2025-01-19 10:25:00+08:00

### 🔴 离线设备 (1)
- **设备**: Smart TV
  - **MAC**: 77:88:99:aa:bb:cc
  - **IP**: 192.168.1.80
  - **厂商**: Samsung Electronics
  - **时间**: 2025-01-19 10:20:00+08:00

### 📋 事件 (2)
- **MAC**: aa:bb:cc:dd:ee:ff | **IP**: 192.168.1.100
  - **类型**: Connected
  - **时间**: 2025-01-19 10:30:00

- **MAC**: 77:88:99:aa:bb:cc | **IP**: 192.168.1.80
  - **类型**: Disconnected
  - **时间**: 2025-01-19 10:20:00

### 🔌 插件事件 (1)
- **插件**: ARPSCAN
  - **MAC**: aa:bb:cc:dd:ee:ff
  - **IP**: 192.168.1.100
  - **状态**: new
```

The notification appears on your iOS device with:
- **Badge count**: 1
- **Group**: NetAlertX (all notifications grouped together)
- **Icon**: Custom icon if `BARK_ICON` is configured

### Example: Test Webhook

```bash
curl -X POST http://localhost:8000/api/webhook/netalertx \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Test notification",
    "username": "NetAlertX",
    "attachments": [...]
  }'
```

## Development

### Project Structure

```
notification-bridge/
├── cmd/server/           # Application entrypoint
├── internal/
│   ├── client/          # External API clients (Bark, etc.)
│   ├── handler/         # HTTP handlers
│   │   └── webhook/     # Webhook-specific handlers
│   ├── middleware/      # HTTP middleware
│   ├── response/        # Response helpers
│   └── types/           # Data type definitions
├── docs/                # Generated Swagger documentation
└── .env                 # Environment configuration (gitignored)
```

### API Documentation

Swagger UI is available at: `http://localhost:8000/swagger/`

### Regenerate Swagger Docs

When you modify API annotations, regenerate the docs:

```bash
swag init -g ./cmd/server/main.go -o ./docs
```

### Architecture Highlights

- **No Framework**: Uses only Go's standard library `net/http`
- **Dependency Injection**: Handler factory pattern for clean dependencies
- **Async Forwarding**: Notifications sent via goroutines (non-blocking)
- **Request Tracking**: Every request gets a unique ID for debugging
- **Type Safety**: Strongly-typed webhook and API structures

### Adding a New Webhook Source

1. Define types in `internal/types/[source].go`
2. Create handler in `internal/handler/webhook/[source].go`
3. Use factory pattern: `func New[Source]Handler(clients...) http.HandlerFunc`
4. Register route in `cmd/server/main.go`
5. Update Swagger docs: `swag init -g ./cmd/server/main.go -o ./docs`

## License

MIT

## Links

- [Bark](https://github.com/Finb/Bark) - iOS push notification service
- [NetAlertX](https://github.com/jokob-sk/NetAlertX) - Network monitoring tool
