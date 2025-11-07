# MyLittlePrice Backend - Elysia.js Edition

Modern, high-performance backend rewrite using Elysia.js and Bun runtime.

## Features

- ⚡ **Elysia.js** - Fast, modern web framework
- 🚀 **Bun Runtime** - Ultra-fast JavaScript runtime
- 🤖 **Google Gemini AI** - Advanced natural language processing with better JS SDK
- 🔍 **SerpAPI** - Product search via Google Shopping
- 🔐 **Google OAuth** - Secure authentication
- 📦 **Redis** - Session and cache management
- 🗄️ **PostgreSQL** - Persistent data storage
- 🎯 **Intelligent Grounding** - Smart decision-making for search grounding
- 🔄 **Key Rotation** - Automatic API key rotation with usage tracking

## Prerequisites

- [Bun](https://bun.sh/) >= 1.0
- Redis >= 6.0
- PostgreSQL >= 12
- Gemini API key(s)
- SerpAPI key(s)
- Google OAuth credentials

## Quick Start

### 1. Install Dependencies

```bash
cd backend-elysia
bun install
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your API keys and configuration
```

### 3. Run Development Server

```bash
bun run dev
```

### 4. Run Production Server

```bash
bun run start
```

## Docker

```bash
# Build image
docker build -t mylittleprice-backend .

# Run container
docker run -p 8080:8080 --env-file .env mylittleprice-backend
```

## API Endpoints

### Chat
- `POST /api/chat` - Send chat message
- `POST /api/product-details` - Get product details

### Authentication
- `GET /api/auth/google/url` - Get Google OAuth URL
- `GET /api/auth/google/callback` - OAuth callback
- `POST /api/auth/refresh` - Refresh tokens
- `GET /api/auth/verify` - Verify token
- `POST /api/auth/logout` - Logout

### Statistics
- `GET /api/health` - Health check
- `GET /api/stats/keys` - API key statistics
- `GET /api/stats/all` - All statistics

## Project Structure

```
backend-elysia/
├── src/
│   ├── config/           # Configuration management
│   ├── modules/          # Feature modules (controllers)
│   │   ├── auth/         # Authentication endpoints
│   │   ├── chat/         # Chat endpoints
│   │   └── stats/        # Statistics endpoints
│   ├── services/         # Business logic
│   │   ├── auth.service.ts
│   │   ├── cache.service.ts
│   │   ├── embedding.service.ts
│   │   ├── gemini.service.ts
│   │   ├── grounding-strategy.service.ts
│   │   ├── search-history.service.ts
│   │   ├── serp.service.ts
│   │   └── session.service.ts
│   ├── utils/            # Utilities
│   │   ├── database.ts
│   │   ├── jwt.ts
│   │   ├── key-rotator.ts
│   │   └── math.ts
│   ├── prompts/          # AI prompts
│   ├── types/            # TypeScript types
│   ├── container.ts      # Dependency injection
│   └── main.ts           # Entry point
├── Dockerfile
├── package.json
└── tsconfig.json
```

## Advantages Over Go Backend

1. **Better Gemini SDK** - @google/generative-ai provides more features and better type safety
2. **Faster Development** - TypeScript and Elysia provide excellent DX
3. **Performance** - Bun runtime is extremely fast
4. **Modern Patterns** - Clean architecture with dependency injection
5. **Type Safety** - Full TypeScript support throughout the codebase

## Configuration

See `.env.example` for all available configuration options.

Key settings:
- `GEMINI_GROUNDING_MODE`: `conservative`, `balanced`, or `aggressive`
- `MAX_SEARCHES_PER_SESSION`: Limit searches per session
- `CORS_ORIGINS`: Comma-separated allowed origins

## Development

```bash
# Run in development mode with hot reload
bun run dev

# Type check
bun run build

# Format code
bun run format
```

## License

MIT
