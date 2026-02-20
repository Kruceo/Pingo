# Pinger - CLI Smoke Ping Service

A CLI service for connectivity monitoring (smoke ping) with metrics storage in ClickHouse.

## 🚀 Features

- ✅ IPv4 and IPv6 ping support
- ✅ JSON configuration
- ✅ ClickHouse storage
- ✅ Descriptive names for each target
- ✅ Configurable ping interval
- ✅ Graceful shutdown
- ✅ Environment variables support

## 📦 Installation

### Prerequisites

- Go 1.18+
- ClickHouse (optional for development)

### Build

#### Normal build:
```bash
go build -o pinger .
```

#### Build with CGO enabled (recommended for better network performance):
```bash
CGO_ENABLED=1 go build -o pingo .
```

#### Build for Linux:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o pingo-linux .
```

#### Optimized build for Linux:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o pingo-linux .
```

## 🛠️ Usage

### 1. Configuration

Create a `config.json` file based on the example:

```bash
cp config.example.json config.json
```

Edit `config.json` with your targets:

```json
{
  "ping_interval": 30,
  "items": [
    {
      "name": "Cloudflare DNS IPv4",
      "tool": "pingv4",
      "target": "1.1.1.1",
      "timeout": 5000
    },
    {
      "name": "Google DNS IPv6",
      "tool": "pingv6",
      "target": "2606:4700:4700::1111",
      "timeout": 5000
    }
  ]
}
```

### 2. ClickHouse Configuration (Optional)

Create a `.env` file to configure ClickHouse:

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
CLICKHOUSE_DSN=localhost:9000
CLICKHOUSE_DATABASE=default
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=yourpassword
PING_INTERVAL=30s
```

### 3. Execution

```bash
# Run the service
./pinger start config.json

# Or with CGO enabled
CGO_ENABLED=1 ./pinger start config.json
```

### 4. Stop the Service

Use `Ctrl+C` for graceful shutdown.

## 📊 Data Structure

### ClickHouse Table

```sql
CREATE TABLE ping_metrics (
    name String,
    target String,
    success UInt8,
    duration_ms Float64,
    error String,
    timestamp DateTime
) ENGINE = MergeTree()
ORDER BY timestamp
```

### Stored Fields

- `name`: Descriptive target name
- `target`: IP or hostname
- `success`: 1 for success, 0 for failure
- `duration_ms`: Response time in milliseconds
- `error`: Error message if failure occurred
- `timestamp`: Measurement date and time

## 🔧 Environment Variables

| Variable | Default | Description |
|-----------|---------|-------------|
| `CLICKHOUSE_DSN` | `localhost:9000` | ClickHouse address |
| `CLICKHOUSE_DATABASE` | `default` | Database name |
| `CLICKHOUSE_USERNAME` | `default` | Username |
| `CLICKHOUSE_PASSWORD` | `` | Password |
| `PING_INTERVAL` | `30s` | Ping interval |

## 🐳 Docker (Optional)

### Run ClickHouse with Docker:

```bash
docker run -d \
  --name clickhouse \
  -p 9000:9000 \
  -p 8123:8123 \
  -v clickhouse_data:/var/lib/clickhouse \
  clickhouse/clickhouse-server
```

### Docker build for the application:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o pinger .

FROM alpine:latest
COPY --from=builder /app/pinger /app/pinger
ENTRYPOINT ["/app/pinger"]
```

## 📋 Build Commands

### Development:
```bash
# Normal build
go build -o pinger .

# Build with CGO
CGO_ENABLED=1 go build -o pinger .
```

### Production Linux:
```bash
# Build for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o pinger-linux .

# Optimized build for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o pinger-linux .
```

### Windows:
```bash
GOOS=windows GOARCH=amd64 go build -o pinger.exe .
```

## 🚨 Troubleshooting

### Raw socket permission error:

```bash
# Linux: run with sudo or set capabilities
sudo setcap cap_net_raw+ep ./pinger

# Or run with sudo
sudo ./pinger start config.json
```

### CGO not available:

Install development dependencies:

```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"

# macOS (with Homebrew)
brew install gcc
```

## 📝 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

1. Fork the project
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request