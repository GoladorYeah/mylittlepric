# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MyLittlePrice is an AI-powered product search assistant that helps users find products through conversational chat. The system uses Google Gemini for natural language processing, SerpAPI for product search via Google Shopping, and maintains session state through Redis. The application consists of a Go backend (Fiber framework) and Next.js frontend with WebSocket real-time communication.

## Architecture

### Backend (Go)
- **Framework**: Fiber v2 (Express-like HTTP framework)
- **Entry Point**: `backend/cmd/api/main.go`
- **Dependency Injection**: Centralized container pattern in `backend/internal/container/container.go`
- **API Communication**:
  - REST endpoints at `/api/*` for HTTP requests
  - WebSocket endpoint at `/ws` for real-time chat
- **Key Services**:
  - `GeminiService`: Handles AI conversation processing with intelligent grounding strategy
  - `SerpService`: Manages product searches via Google Shopping API with relevance scoring
  - `SessionService`: Manages chat sessions with Redis persistence
  - `CacheService`: Caches search results and embeddings
  - `EmbeddingService`: Category detection using semantic embeddings
  - `GroundingStrategy`: Smart decision-making for when to use Google Search grounding

### Frontend (Next.js)
- **Framework**: Next.js 15 (App Router with Turbopack)
- **State Management**: Zustand with localStorage persistence
- **Styling**: Tailwind CSS v4 with dark mode support (next-themes)
- **Real-time**: WebSocket via react-use-websocket
- **Main Pages**:
  - `/`: Landing page
  - `/chat`: Main chat interface (can accept `?q=` query param)

### Data Flow
1. User sends message → WebSocket (`/ws`) or REST (`/api/chat`)
2. Session validated/created in Redis
3. Message sent to GeminiService with conversation history
4. Gemini returns either:
   - `response_type: "dialogue"` → Display text + quick replies
   - `response_type: "search"` → Trigger SerpService product search
5. If search: Query translated to English → SerpAPI → Results scored for relevance
6. Products + response sent back to frontend
7. Session state updated (category, last product, search count)

### Key Rotator Pattern
The system uses API key rotation to handle rate limits:
- `utils.KeyRotator` manages multiple API keys for both Gemini and SerpAPI
- Automatically rotates to next key on quota errors
- Tracks usage statistics per key in Redis
- View stats at `/api/stats/keys`

### Intelligent Grounding Strategy
The system dynamically decides when to enable Google Search grounding:
- **Conservative mode**: Grounding only for product-specific queries
- **Balanced mode** (default): Adaptive based on query complexity
- **Aggressive mode**: Grounding for most queries
- Tracks decision statistics at `/api/stats/grounding`

## Development Commands

### Backend (Go)
```bash
# Navigate to backend
cd backend

# Install dependencies
go mod download

# Run development server (requires .env file)
go run cmd/api/main.go

# Build binary
go build -o bin/api cmd/api/main.go
```

### Frontend (Next.js)
```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Run development server (with Turbopack)
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

### Docker Compose
```bash
# Start all services (PostgreSQL, Redis, Backend, Frontend)
docker-compose up

# Start specific service
docker-compose up backend

# Rebuild and start
docker-compose up --build

# Stop all services
docker-compose down

# View logs
docker-compose logs -f backend
```

## Configuration

### Backend Environment Variables (.env)
Required variables (see `backend/internal/config/config.go` for full list):
- `GEMINI_API_KEYS`: Comma-separated Gemini API keys
- `SERP_API_KEYS`: Comma-separated SerpAPI keys
- `REDIS_URL`: Redis connection string (default: localhost:6379)
- `PORT`: Server port (default: 8080)
- `CORS_ORIGINS`: Allowed origins (default: http://localhost:3000)

Grounding configuration:
- `GEMINI_USE_GROUNDING`: Enable/disable grounding (default: true)
- `GEMINI_GROUNDING_MODE`: "conservative", "balanced", or "aggressive" (default: balanced)
- `GEMINI_GROUNDING_MIN_WORDS`: Minimum query words for grounding (default: 2)

Cache TTL (seconds):
- `CACHE_SERP_TTL`: SerpAPI results cache (default: 86400 = 24h)
- `CACHE_GEMINI_TTL`: Gemini responses cache (default: 3600 = 1h)
- `CACHE_IMMERSIVE_TTL`: Product details cache (default: 43200 = 12h)

Session limits:
- `SESSION_TTL`: Session expiration (default: 86400 = 24h)
- `MAX_MESSAGES_PER_SESSION`: Max messages (default: 8)
- `MAX_SEARCHES_PER_SESSION`: Max product searches (default: 3)

### Frontend Environment Variables
- `NEXT_PUBLIC_API_URL`: Backend API URL (default: http://localhost:8080)

## Important Implementation Details

### Prompt Management
The system uses specialized prompts based on detected product category:
- `master_prompt.txt`: Default conversational prompt
- `specialized_electronics.txt`: For electronics/tech products
- `specialized_parametric.txt`: For products with technical specs
- `specialized_generic_model.txt`: For general product models

Category detection uses embeddings (`EmbeddingService`) to match user queries to predefined categories.

### Search Relevance Scoring
`SerpService.calculateRelevanceScore()` implements sophisticated product matching:
- Full phrase matching (highest score)
- Word order preservation
- Brand detection
- Model number matching
- Filters common words
- Adjustable thresholds per search type (exact/parameters/category)

### Session State Management
Sessions track:
- Conversation history (limited by MAX_MESSAGES_PER_SESSION)
- Current product category
- Last selected product (price/name)
- Search count (limited by MAX_SEARCHES_PER_SESSION)
- Search status (idle/searching/completed)

### WebSocket vs REST
Both handlers (`WSHandler` and `ChatHandler`) share similar logic but:
- WebSocket: Real-time bidirectional communication, automatic reconnection
- REST: Simpler for single request/response, easier debugging
- Both support session management and search functionality

### Translation for Better Results
Search queries are automatically translated to English before sending to SerpAPI to improve result quality across different markets (see `GeminiService.TranslateToEnglish()`).

## API Endpoints

### Chat
- `POST /api/chat`: HTTP chat endpoint
- `GET /ws`: WebSocket chat endpoint

### Product Details
- `POST /api/product-details`: Get detailed product info by page token

### Statistics & Monitoring
- `GET /api/stats/keys`: API key usage statistics
- `GET /api/stats/grounding`: Grounding decision statistics
- `GET /api/stats/tokens`: Token usage metrics
- `GET /api/stats/all`: Combined statistics

### Health Check
- `GET /health`: Service health status

## Testing

### Manual API Testing
Use the provided PowerShell script:
```powershell
.\test-api.ps1
```

### Product Search Flow Testing
1. Start chat: "I need a laptop"
2. Assistant asks clarifying questions
3. Provide specifications: "gaming laptop under $1500"
4. System searches and returns products
5. Ask follow-up: "show me cheaper options"

## Common Patterns

### Adding a New Service
1. Create service in `backend/internal/services/`
2. Initialize in `container.initServices()`
3. Inject container into handlers
4. Access via `h.container.YourService`

### Modifying Conversation Logic
- Edit prompts in `backend/internal/services/prompts/`
- Adjust grounding strategy in `backend/internal/services/grounding_strategy.go`
- Update response parsing in `backend/internal/handlers/chat.go` or `websoket.go`

### Adding Frontend State
- Add to `ChatStore` interface in `frontend/src/lib/store.ts`
- Implement getter/setter methods
- Use in components via `useChatStore()`

### Debugging Search Results
Enable verbose logging by checking console output:
- Search query and parameters logged in `SerpService.SearchProducts()`
- Top 5 results with relevance scores displayed
- Translation logs show original → English query conversion

## Database Schema

PostgreSQL (via migrations in `backend/migrations/`):
- Currently minimal usage - Redis handles most session state
- Schema available in `001_initial_schema.sql`

Redis Keys:
- `session:{session_id}`: Serialized session data
- `messages:{session_id}`: Conversation history
- `cache:search:{country}:{type}:{query}`: Cached search results
- `cache:product:{page_token}`: Cached product details
- `keyrotator:{service}:*`: API key rotation state
- `embeddings:categories`: Category embeddings cache
