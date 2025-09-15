# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

New API is a next-generation AI gateway and asset management system built on Go and React. It's based on One API and provides comprehensive management of AI models, channels, and user interactions with support for multiple AI providers.

## Development Commands

### Backend (Go)
- `go run main.go` - Start the backend development server
- `make start-backend` - Start backend using Makefile (runs Go app in background)
- `make all` - Build frontend and start backend

### Frontend (React)
- `cd web && bun install` - Install frontend dependencies
- `cd web && bun run dev` - Start frontend development server
- `cd web && bun run build` - Build frontend for production
- `cd web && bun run lint` - Check code formatting
- `cd web && bun run lint:fix` - Fix code formatting issues
- `make build-frontend` - Build frontend using Makefile

### Docker Development
- `docker-compose up -d` - Start full stack with MySQL and Redis
- `docker run --name new-api -d --restart always -p 3000:3000 -e TZ=Asia/Shanghai -v /home/ubuntu/data/new-api:/data calciumion/new-api:latest` - Run with SQLite

## Architecture Overview

### Backend Structure (Go)
- **main.go** - Application entry point with embedded frontend assets
- **controller/** - HTTP request handlers for API endpoints
- **middleware/** - Request middleware (auth, rate limiting, CORS, etc.)
- **model/** - Database models and business logic
- **common/** - Shared utilities and helper functions
- **constant/** - Application constants and enums
- **router/** - Route definitions and setup
- **service/** - Business services and external integrations
- **dto/** - Data transfer objects for API requests/responses

### Frontend Structure (React)
- Built with React 18, Semi-UI components, and Vite
- **src/components/** - Organized by feature (dashboard, auth, settings, etc.)
- **src/components/layout/** - Layout components and navigation
- **src/components/table/** - Data table components for various entities
- **src/components/playground/** - AI model testing interface

### Key Features
- Multi-model AI gateway supporting OpenAI, Claude, Gemini, and 20+ providers
- User management with quota tracking and billing
- Channel management with load balancing and failover
- Real-time API monitoring and usage analytics
- Web-based playground for testing AI models
- Support for various AI tasks: chat, image generation, rerank, etc.

### Database
- Supports SQLite (default), MySQL 5.7.8+, and PostgreSQL 9.6+
- Uses GORM for database operations
- Embedded SQLite for development, external DB for production

### Environment Configuration
Key environment variables:
- `SQL_DSN` - Database connection string
- `REDIS_CONN_STRING` - Redis connection for caching
- `SESSION_SECRET` - Required for multi-node deployments
- `GIN_MODE=debug` - Enable debug mode
- `PORT` - Server port (default 3000)

### Multi-Node Deployment
- Supports horizontal scaling with Redis for session storage
- Set `SESSION_SECRET` and `CRYPTO_SECRET` for multi-node setups
- Use `NODE_TYPE=slave` for worker nodes

## Development Notes

- Frontend uses bun as package manager instead of npm/yarn
- Backend embeds frontend assets using Go embed
- Hot reloading available for both frontend (Vite) and backend (go run)
- The application serves both API endpoints and the React SPA from the same port
- Default development setup runs on port 3000