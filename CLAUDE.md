# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MyLittlePrice is a full-stack AI-powered price comparison application that helps users find and compare product prices. The application uses Gemini AI for natural language processing and Google SERP API for product searches.

**Architecture**: Monorepo with separate frontend and backend applications
- **Frontend**: Next.js 16 (React 19) with TypeScript, Tailwind CSS v4
- **Backend**: Go (Fiber framework) with PostgreSQL and Redis
- **Infrastructure**: Docker Compose for local development, Grafana/Loki for monitoring

## Development Commands

### Backend (Go)
```bash
cd backend

# Run in development mode
go run cmd/api/main.go

# Build executable
go build -o api.exe cmd/api/main.go

# Run the built executable
./api.exe

# Generate Ent models after schema changes
go generate ./ent

# Install dependencies
go mod download
go mod tidy
```

### Frontend (Next.js)
```bash
cd frontend

# Development server with Turbopack
npm run dev

# Production build
npm run build

# Start production server
npm run start

# Install dependencies
npm install
```

### Infrastructure

```bash
# Start database services (PostgreSQL + Redis)
docker-compose up -d

# Start monitoring stack (Grafana + Loki + Promtail)
docker-compose -f docker-compose.monitoring.yml up -d

# Stop all services
docker-compose down
docker-compose -f docker-compose.monitoring.yml down

# View logs
docker-compose logs -f postgres
docker-compose logs -f redis
```

**Service Ports**:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Grafana: http://localhost:3001 (admin/admin)
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- Loki: localhost:3100

## Environment Configuration

Both frontend and backend require `.env` files. Examples are provided in `.env.example` files.

**Critical Backend Configuration**:
- `GEMINI_API_KEYS`: Comma-separated Gemini API keys for AI features
- `SERP_API_KEYS`: Comma-separated SERP API keys for product search
- `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET`: OAuth credentials (must match frontend)
- `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET`: Generate with `openssl rand -hex 32`
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `CORS_ORIGINS`: Comma-separated frontend URLs (NO trailing slashes)

**Critical Frontend Configuration**:
- `NEXT_PUBLIC_GOOGLE_CLIENT_ID`: Must match backend OAuth Client ID
- `NEXT_PUBLIC_API_URL`: Backend API URL (default: http://localhost:8080)

## Architecture

### Backend Architecture (Go)

**Framework**: Go Fiber v2 (Express-like HTTP framework)

**Key Architectural Patterns**:
1. **Dependency Injection via Container**: [internal/container/container.go](internal/container/container.go) - Central IoC container that initializes and manages all services, database connections, and API clients
2. **Service Layer Pattern**: Business logic in [internal/services/](internal/services/)
3. **Handler Layer**: HTTP handlers in [internal/handlers/](internal/handlers/)
4. **Domain Models**: Domain types in [internal/domain/](internal/domain/)
5. **Ent ORM**: Database schema defined in [backend/ent/schema/](backend/ent/schema/) - auto-generates type-safe query builders

**Entry Point**: [backend/cmd/api/main.go](backend/cmd/api/main.go)
- Initializes container with all dependencies
- Sets up Fiber app with middleware (CORS, logger, recovery, auth)
- Configures routes via [internal/app/routes.go](internal/app/routes.go)
- Starts cleanup job for expired sessions
- Handles graceful shutdown

**Core Services** (all initialized in container):
- `GeminiService`: AI chat responses with smart grounding (decides when to use Google Search)
- `SerpService`: Product searches via Google SERP API with relevance scoring
- `SessionService`: Manages chat sessions (Redis + PostgreSQL sync)
- `MessageService`: Stores conversation history in Redis
- `AuthService`: JWT authentication + Google OAuth
- `EmbeddingService`: Text embeddings for semantic similarity (category detection, query comparison)
- `CacheService`: Redis caching for API responses (Gemini, SERP, embeddings)
- `SearchHistoryService`: User search history with product click tracking
- `PreferencesService`: User preferences (country, language)

**Key Features**:
- **API Key Rotation**: [internal/utils/key_rotator.go](internal/utils/key_rotator.go) - Round-robin rotation with Redis-based error tracking
- **Smart Grounding Strategy**: [internal/services/grounding_strategy.go](internal/services/grounding_strategy.go) - Intelligent decision-making for when to use Google Search grounding based on query type, freshness needs, and dialogue context
- **WebSocket Support**: Real-time chat via [internal/handlers/websocket.go](internal/handlers/websocket.go)
- **Optional Authentication**: Most endpoints work for both authenticated and anonymous users (using session IDs)

**Database**:
- **ORM**: Ent (Facebook's entity framework) with schema-first approach
- **Migrations**: SQL files in [backend/migrations/](backend/migrations/) - run automatically on container startup
- **Models**: `User`, `ChatSession`, `Message`, `SearchHistory`, `UserPreference`

### Frontend Architecture (Next.js)

**Framework**: Next.js 16 with App Router, React 19, TypeScript

**Project Structure** (Feature-Sliced Design inspired):
```
src/
├── app/                    # Next.js App Router
│   ├── (app)/             # Main app pages (chat, product search)
│   ├── (auth)/            # Authentication pages (login, signup)
│   ├── (marketing)/       # Marketing pages (landing, about)
│   ├── layout.tsx         # Root layout with theme provider
│   └── globals.css        # Global styles + Tailwind
├── features/              # Feature modules
│   ├── auth/             # Authentication logic, hooks, components
│   ├── chat/             # Chat interface and WebSocket logic
│   ├── search/           # Product search UI
│   ├── products/         # Product display components
│   └── policies/         # Privacy/terms pages
└── shared/               # Shared utilities
    ├── components/       # Reusable UI components
    ├── hooks/            # Custom React hooks (API, localStorage, etc.)
    ├── lib/              # Utilities and helpers
    └── types/            # TypeScript type definitions
```

**State Management**:
- Zustand for global state ([features/auth/store.ts](features/auth/store.ts) for auth state)
- React hooks for local state
- WebSocket connections via `react-use-websocket`

**Styling**:
- Tailwind CSS v4 (latest version)
- CSS variables for theming in [globals.css](frontend/src/app/globals.css)
- `next-themes` for dark mode support
- Utility functions: `cn()` from `tailwind-merge` + `clsx`

**Path Aliases** (see [tsconfig.json](frontend/tsconfig.json)):
- `@/*` → `src/*`
- `@/features/*` → `src/features/*`
- `@/shared/*` → `src/shared/*`

**API Communication**:
- REST API calls via custom `useApi` hook
- WebSocket for real-time chat
- Next.js rewrites proxy `/api/*` and `/ws` to backend (configured in [next.config.ts](frontend/next.config.ts))

### Communication Flow

1. **Chat Interaction**: User sends message → WebSocket → Backend processes with Gemini AI → May trigger SERP search → Response streamed back
2. **Product Search**: User query → Backend analyzes with embeddings → SERP API search → Results filtered by relevance → Cached in Redis
3. **Authentication**: Google OAuth popup → Backend validates token → JWT issued → Stored in localStorage + Zustand

### Smart Grounding System

The backend uses an intelligent "Smart Grounding" system to decide when to use Google Search grounding with Gemini:

**Modes** (configured via `GEMINI_GROUNDING_MODE`):
- `conservative`: Minimal usage (~10-20% of requests) - cheapest
- `balanced`: Optimal balance (~30-40% of requests) - **RECOMMENDED**
- `aggressive`: Maximum usage (~60-80% of requests) - most accurate

**Decision Factors** (see [grounding_strategy.go](backend/internal/services/grounding_strategy.go)):
- Fresh info needs (recent products, current prices)
- Specific product detection (brand + model number)
- Dialogue drift (conversation deviates from initial search)
- Electronics category (benefits more from fresh data)

### Key Technical Decisions

1. **Monorepo Structure**: Frontend and backend in same repo for easier coordination
2. **Optional Auth**: Most features work without login - sessions identified by anonymous session IDs
3. **Ent ORM**: Type-safe, schema-first ORM with auto-generated query builders
4. **Redis for Sessions**: Fast session/message storage with PostgreSQL sync for persistence
5. **API Key Rotation**: Multiple API keys with automatic rotation on errors/rate limits
6. **Embedding-based Filtering**: Semantic similarity for category detection and query comparison
7. **Next.js Turbopack**: Faster development builds (still experimental)

## Common Development Workflows

### Adding a New Backend Route

1. Define handler in [internal/handlers/](internal/handlers/)
2. Add route in [internal/app/routes.go](internal/app/routes.go)
3. Add middleware if needed (auth, rate limiting)

### Adding New Database Model (Ent)

1. Create schema in [backend/ent/schema/](backend/ent/schema/)
2. Run `go generate ./ent` to generate code
3. Create migration SQL in [backend/migrations/](backend/migrations/)
4. Restart backend to apply migration

### Adding New Service

1. Create service in [internal/services/](internal/services/)
2. Add to container in [internal/container/container.go](internal/container/container.go)
3. Initialize in `initServices()` method
4. Inject into handlers as needed

### Adding Frontend Feature

1. Create feature directory in [src/features/](src/features/)
2. Structure: `components/`, `hooks/`, `index.ts`
3. Add page in [src/app/](src/app/) if needed
4. Use path aliases for imports

## Important Notes

- **CORS Configuration**: Backend `CORS_ORIGINS` must match frontend URLs exactly (no trailing slashes)
- **OAuth Setup**: Google OAuth Client ID must be identical in both frontend and backend `.env` files
- **Database Migrations**: Auto-run on startup - SQL files executed in order
- **API Keys**: Use comma-separated values for multiple keys (automatic rotation)
- **WebSocket Auth**: Token passed via query parameter (`?token=...`) for WebSocket connections
- **Redis TTLs**: Different cache durations for different data types (see [backend/.env.example](backend/.env.example))

## Monitoring & Debugging

- **Logs**: Check Grafana dashboard at http://localhost:3001 after starting monitoring stack
- **Health Check**: GET http://localhost:8080/health for backend status
- **Stats Endpoints**:
  - `/api/stats/keys` - API key rotation status
  - `/api/stats/grounding` - Smart grounding statistics
  - `/api/stats/tokens` - Token usage tracking
  - `/api/stats/all` - All stats combined
