# Build and Deployment Guide

This guide covers how to build, compile, and restart the New-API application services.

## Prerequisites

- Go 1.19 or later
- Node.js and npm/bun (for frontend)
- Git

## Project Structure

```
new-api/
├── web/                 # Frontend React application
├── relay/              # API relay logic
├── setting/            # Configuration settings
├── service/            # Core services
├── main.go            # Application entry point
├── Dockerfile         # Docker build configuration
└── docs/              # Documentation
```

## Frontend Development

### Installing Dependencies

```bash
cd web
npm install
# OR
bun install
```

### Development Server

Start the frontend development server:

```bash
cd web
npm run dev
# OR
bun run dev
```

The development server will run on `http://localhost:3001` by default.

### Building for Production

```bash
cd web
npm run build
# OR
DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$(cat ../VERSION) bun run build
```

This creates optimized production files in `web/dist/`.

## Backend Development

### Building the Go Application

From the project root directory:

```bash
# Basic build
go build -o new-api

# Production build with version and optimizations
go build -ldflags "-s -w -X 'one-api/common.Version=$(cat VERSION)'" -o new-api
```

### Build Flags Explanation:
- `-s`: Strip symbol table and debug info
- `-w`: Strip DWARF debug info
- `-X 'one-api/common.Version=$(cat VERSION)'`: Inject version from VERSION file

### Running the Application

```bash
# Run directly
./new-api

# Run in background
./new-api &

# Run with custom environment
PORT=3000 ./new-api
```

## Service Management

### Checking Running Services

```bash
# Check for new-api processes
ps aux | grep new-api

# Check port usage
lsof -i :3000
lsof -i :3001
```

### Stopping Services

```bash
# Stop by process name
pkill new-api

# Stop by port (force kill)
lsof -ti:3000 | xargs kill -9
lsof -ti:3001 | xargs kill -9

# Stop frontend development servers
pkill -f vite
pkill -f esbuild
```

### Complete Service Restart

```bash
#!/bin/bash
# Stop all related processes
pkill new-api
pkill -f vite
pkill -f esbuild
lsof -ti:3000 | xargs kill -9 2>/dev/null
lsof -ti:3001 | xargs kill -9 2>/dev/null

# Rebuild the application
go build -ldflags "-s -w -X 'one-api/common.Version=$(cat VERSION)'" -o new-api

# Start the service
./new-api &

echo "Service restarted successfully"
```

## Environment Configuration

### Environment Variables

Create a `.env` file in the project root:

```bash
# Database
SQL_DSN="sqlite:///data/one-api.db"
# SQL_DSN="mysql://username:password@tcp(localhost:3306)/oneapi"
# SQL_DSN="postgres://username:password@localhost:5432/oneapi"

# Redis (optional)
REDIS_CONN_STRING="redis://localhost:6379"

# Server
PORT=3000
HOST=0.0.0.0

# Logging
LOG_LEVEL=info

# Security
JWT_SECRET="your-secret-key"
SESSION_SECRET="your-session-secret"
```

### Default Configuration

If no `.env` file is found, the application uses these defaults:
- Database: SQLite (`./data/one-api.db`)
- Port: 3000
- Host: 0.0.0.0
- No Redis caching

## Docker Deployment

### Building Docker Image

```bash
# Build image
docker build -t new-api:latest .

# Build with custom tag
docker build -t new-api:$(cat VERSION) .
```

### Running with Docker

```bash
# Run container
docker run -d \
  --name new-api \
  -p 3000:3000 \
  -v $(pwd)/data:/data \
  new-api:latest

# Run with environment file
docker run -d \
  --name new-api \
  -p 3000:3000 \
  -v $(pwd)/data:/data \
  --env-file .env \
  new-api:latest
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  new-api:
    build: .
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./.env:/app/.env
    restart: unless-stopped
    environment:
      - SQL_DSN=sqlite:///data/one-api.db
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/api/status"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Run with:
```bash
docker-compose up -d
```

## Development Workflow

### 1. Making Code Changes

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make your changes...

# Test locally
go build -o new-api
./new-api &
```

### 2. Testing Changes

```bash
# Run tests (if available)
go test ./...

# Manual testing
curl http://localhost:3000/api/status
```

### 3. Building for Production

```bash
# Clean previous builds
rm -f new-api

# Build optimized binary
go build -ldflags "-s -w -X 'one-api/common.Version=$(cat VERSION)'" -o new-api

# Test production build
./new-api &
```

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Find and kill process using port 3000
lsof -ti:3000 | xargs kill -9
```

#### Build Failures
```bash
# Clean module cache
go clean -modcache
go mod download
go mod tidy
```

#### Frontend Build Issues
```bash
cd web
rm -rf node_modules package-lock.json
npm install
npm run build
```

### Log Analysis

```bash
# View recent logs
tail -f logs/app.log

# Check system logs
journalctl -u new-api -f

# Debug mode
DEBUG=true ./new-api
```

### Performance Monitoring

```bash
# Memory usage
ps aux | grep new-api

# CPU and memory details
top -p $(pgrep new-api)

# Network connections
netstat -an | grep :3000
```

## Production Deployment Checklist

- [ ] Environment variables configured
- [ ] Database connection tested
- [ ] SSL/TLS certificates configured
- [ ] Firewall rules set up
- [ ] Log rotation configured
- [ ] Backup strategy in place
- [ ] Monitoring and alerting set up
- [ ] Health checks configured
- [ ] Resource limits set
- [ ] Security headers configured

## Monitoring and Maintenance

### Health Checks

```bash
# API health check
curl -f http://localhost:3000/api/status

# Database health
curl -f http://localhost:3000/api/health/db

# System metrics
curl -f http://localhost:3000/api/metrics
```

### Regular Maintenance

```bash
# Update dependencies
go mod tidy
go mod download

# Clean old builds
rm -f old-binaries/new-api-*

# Rotate logs
logrotate /etc/logrotate.d/new-api
```

This guide should help you efficiently build, deploy, and maintain the New-API application in various environments.