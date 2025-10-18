# Architecture Overview

Visual representation of MyLittlePrice system architecture.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           User Browser                               │
│                     (mylittleprice.com)                             │
└───────────────────────────┬─────────────────────────────────────────┘
                            │
                            │ HTTPS (Cloudflare Tunnel)
                            ↓
┌─────────────────────────────────────────────────────────────────────┐
│                      Cloudflare CDN/Tunnel                          │
│  • SSL/TLS Termination                                              │
│  • DDoS Protection                                                  │
│  • Caching                                                          │
└───────────┬────────────────────────────────┬────────────────────────┘
            │                                │
            │ :3000                          │ :8080
            ↓                                ↓
┌───────────────────────┐        ┌──────────────────────────┐
│   Frontend Container  │        │   Backend Container      │
│   (Next.js 15)        │◄──────►│   (Go + Fiber)          │
│                       │        │                          │
│  • App Router         │  WS/   │  • REST API             │
│  • Zustand Store      │  HTTP  │  • WebSocket Handler    │
│  • Tailwind CSS       │        │  • Session Management   │
│  • Dark Mode          │        │  • Key Rotation         │
│                       │        │                          │
│  Port: 3000          │        │  Port: 8080             │
└───────────────────────┘        └────────┬─────────────────┘
                                          │
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    │                     │                     │
                    ↓                     ↓                     ↓
         ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
         │ PostgreSQL       │  │ Redis            │  │ External APIs    │
         │ Container        │  │ Container        │  │                  │
         │                  │  │                  │  │ • Google Gemini  │
         │ • User Data      │  │ • Session Cache  │  │ • SerpAPI        │
         │ • Migrations     │  │ • Search Cache   │  │ • Google Shopping│
         │ • Persistence    │  │ • Key Stats      │  │                  │
         │                  │  │ • AOF Enabled    │  │                  │
         │ Port: 5432      │  │ Port: 6379      │  │                  │
         └──────────────────┘  └──────────────────┘  └──────────────────┘
```

## Docker Network Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                      Docker Host Machine                             │
│  ┌────────────────────────────────────────────────────────────────┐  │
│  │            Docker Network: mylittleprice_default               │  │
│  │                    (Bridge Network)                            │  │
│  │                                                                │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌──────────┐  ┌────────┐  │  │
│  │  │  frontend   │  │   backend   │  │ postgres │  │ redis  │  │  │
│  │  │  :3000      │──│   :8080     │──│  :5432   │──│ :6379  │  │  │
│  │  └──────┬──────┘  └──────┬──────┘  └────┬─────┘  └───┬────┘  │  │
│  │         │                │              │            │        │  │
│  └─────────┼────────────────┼──────────────┼────────────┼────────┘  │
│            │                │              │            │           │
│  ┌─────────┼────────────────┼──────────────┼────────────┼────────┐  │
│  │ Volumes │                │              │            │        │  │
│  │         ↓                │              ↓            ↓        │  │
│  │  node_modules       (no vol)    postgres_data   redis_data   │  │
│  │  .next (dev)                                                  │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  Port Mappings (Host:Container)                                     │
│  • 3000:3000  → Frontend                                            │
│  • 8080:8080  → Backend                                             │
│  • 5432:5432  → PostgreSQL                                          │
│  • 6379:6379  → Redis                                               │
└──────────────────────────────────────────────────────────────────────┘
```

## Data Flow - Search Request

```
User Types: "I need a laptop under $1000"
     │
     ↓
┌─────────────────────────────────────────────────────────────────┐
│ 1. Frontend (React Component)                                   │
│    • useChatStore.sendMessage()                                 │
│    • WebSocket or HTTP POST                                     │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. Backend (WebSocket/HTTP Handler)                            │
│    • Validate session                                           │
│    • Check rate limits                                          │
│    • Parse message                                              │
└────────────────────────┬────────────────────────────────────────┘
                         │
        ┌────────────────┴────────────────┐
        │                                 │
        ↓                                 ↓
┌──────────────────┐            ┌─────────────────────┐
│ 3a. SessionService│            │ 3b. CacheService    │
│ • Get/Create      │            │ • Check cache       │
│ • Redis: session: │            │ • Redis: cache:     │
└────────┬──────────┘            └──────────┬──────────┘
         │                                  │
         └───────────┬──────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────────┐
│ 4. GeminiService                                                │
│    • Build conversation context                                 │
│    • Select appropriate prompt (category-based)                 │
│    • Check if grounding needed (GroundingStrategy)              │
│    • Call Gemini API with history                               │
│    • Parse response (dialogue vs search)                        │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ↓
              ┌──────────┴──────────┐
              │                     │
              ↓                     ↓
    ┌─────────────────┐   ┌──────────────────────┐
    │ 5a. Dialogue    │   │ 5b. Search Request   │
    │ • Quick replies │   │ • Extract query      │
    │ • Send to user  │   │ • Translate to EN    │
    └─────────────────┘   └──────────┬───────────┘
                                     │
                                     ↓
                          ┌─────────────────────────┐
                          │ 6. SerpService          │
                          │ • Search Google Shopping│
                          │ • Relevance scoring     │
                          │ • Filter results        │
                          │ • Cache results         │
                          └──────────┬──────────────┘
                                     │
                                     ↓
                          ┌─────────────────────────┐
                          │ 7. Format Response      │
                          │ • Product list          │
                          │ • Metadata              │
                          │ • Quick actions         │
                          └──────────┬──────────────┘
                                     │
                                     ↓
                          ┌─────────────────────────┐
                          │ 8. Update Session       │
                          │ • Increment searches    │
                          │ • Save last product     │
                          │ • Update category       │
                          └──────────┬──────────────┘
                                     │
                                     ↓
                          ┌─────────────────────────┐
                          │ 9. Send to Frontend     │
                          │ • WebSocket or HTTP     │
                          │ • JSON response         │
                          └──────────┬──────────────┘
                                     │
                                     ↓
                          ┌─────────────────────────┐
                          │ 10. Frontend Display    │
                          │ • Update chat messages  │
                          │ • Render products       │
                          │ • Show quick replies    │
                          └─────────────────────────┘
```

## API Key Rotation Flow

```
Backend Service needs API key
        │
        ↓
┌──────────────────────────────┐
│ KeyRotator (Gemini/Serp)     │
│ • Check current key          │
│ • Track usage stats          │
└────────────┬─────────────────┘
             │
             ↓
      ┌──────┴──────┐
      │             │
      ↓             ↓
  Available?     Quota Error?
      │             │
      Yes           Yes
      │             │
      ↓             ↓
  Use Key      Rotate Next
      │             │
      ↓             ↓
  ┌────────────┐ ┌──────────────┐
  │ Update     │ │ Mark Failed  │
  │ Stats in   │ │ Try Next Key │
  │ Redis      │ │              │
  └────────────┘ └──────────────┘
```

## Grounding Strategy Decision Tree

```
User Query Received
        │
        ↓
┌─────────────────────────────────────────┐
│ GroundingStrategy.ShouldUseGrounding()  │
└────────────┬────────────────────────────┘
             │
             ↓
    ┌────────┴────────┐
    │                 │
    ↓                 ↓
Contains        Analyze Features:
product/model   • Dialogue drift
specific?       • Electronics category
    │           • Fresh info needed
    Yes         • Specific product
    │                 │
    ↓                 ↓
  Score: 1.0    Calculate weighted score
    │                 │
    └────────┬────────┘
             ↓
    Score > Threshold (0.5)?
             │
    ┌────────┴────────┐
    │                 │
    Yes               No
    │                 │
    ↓                 ↓
Use Grounding   No Grounding
    │                 │
    └────────┬────────┘
             ↓
    Update Statistics
    (Redis: stats:grounding)
```

## Session Lifecycle

```
User Opens Chat
      │
      ↓
┌─────────────────────────┐
│ No Session ID?          │
│ Generate UUID           │
│ Create in Redis         │
└──────────┬──────────────┘
           │
           ↓
┌─────────────────────────┐
│ Session Active          │
│ • Track messages (max 8)│
│ • Track searches (max 3)│
│ • Save category         │
│ • TTL: 24 hours         │
└──────────┬──────────────┘
           │
    ┌──────┴──────┐
    │             │
    ↓             ↓
Limits       TTL Expired?
Reached?          │
    │             │
    Yes           Yes
    │             │
    ↓             ↓
Block New    Delete Session
Searches     Clear Messages
    │             │
    └──────┬──────┘
           │
           ↓
┌─────────────────────────┐
│ Session Ends            │
│ • User closes tab       │
│ • TTL expires           │
│ • Manual clear          │
└─────────────────────────┘
```

## Caching Strategy

```
Request Comes In
      │
      ↓
┌─────────────────────────────────────────┐
│ Cache Key Generation                    │
│ • Type: search/gemini/product/embedding │
│ • Parameters: country, query, etc.      │
└──────────┬──────────────────────────────┘
           │
           ↓
    ┌──────┴──────┐
    │             │
    ↓             ↓
Cache Hit?    Cache Miss
    │             │
    Yes           No
    │             │
    ↓             ↓
Return        Fetch from API
Cached             │
Data               ↓
    │         ┌────────────┐
    │         │ Store with │
    │         │ TTL        │
    │         └──────┬─────┘
    │                │
    └────────┬───────┘
             ↓
        Return Data

Cache TTLs:
• SerpAPI: 24h
• Gemini: 1h
• Product Details: 12h
• Embeddings: 24h
```

## Service Dependencies

```
┌─────────────────────────────────────────────────────────────────┐
│                    Dependency Container                         │
│  (internal/container/container.go)                              │
│                                                                 │
│  Initialization Order:                                          │
│                                                                 │
│  1. Config       ← Environment variables                        │
│  2. Redis        ← Config.REDIS_URL                             │
│  3. PostgreSQL   ← Config.DATABASE_URL                          │
│  4. CacheService ← Redis                                        │
│  5. KeyRotators  ← Config.API_KEYS                              │
│  6. Embedding    ← Gemini, Cache                                │
│  7. Grounding    ← Config                                       │
│  8. Session      ← Redis, Config                                │
│  9. Gemini       ← KeyRotator, Cache, Embedding, Grounding      │
│ 10. Serp         ← KeyRotator, Cache                            │
│ 11. Handlers     ← All Services                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Deployment Architecture (Production)

```
                    Internet
                       │
                       ↓
            ┌──────────────────────┐
            │  Cloudflare Tunnel   │
            │  • SSL/TLS           │
            │  • DDoS Protection   │
            │  • Rate Limiting     │
            └──────────┬───────────┘
                       │
        ┌──────────────┴──────────────┐
        │                             │
        ↓                             ↓
mylittleprice.com            api.mylittleprice.com
        │                             │
        ↓                             ↓
    localhost:3000              localhost:8080
        │                             │
        └──────────────┬──────────────┘
                       │
         ┌─────────────┴─────────────┐
         │  Docker Host (Server)     │
         │                           │
         │  docker-compose.prod.yml  │
         │                           │
         │  ┌────────────────────┐   │
         │  │ Network Bridge     │   │
         │  │ • frontend         │   │
         │  │ • backend          │   │
         │  │ • postgres (local) │   │
         │  │ • redis (local)    │   │
         │  └────────────────────┘   │
         │                           │
         │  Volumes:                 │
         │  • postgres_data          │
         │  • redis_data             │
         │                           │
         │  Auto-restart: always     │
         │  Health checks: enabled   │
         └───────────────────────────┘
```

## Tech Stack Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                        Frontend Stack                            │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │ Next.js 15 (App Router) + React 19 + TypeScript         │    │
│  │ • Turbopack (dev/prod)                                   │    │
│  │ • Tailwind CSS v4                                        │    │
│  │ • Zustand (state)                                        │    │
│  │ • react-use-websocket                                    │    │
│  │ • next-themes (dark mode)                                │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                        Backend Stack                             │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │ Go 1.25 + Fiber v2                                       │    │
│  │ • google/generative-ai-go (Gemini)                       │    │
│  │ • redis/go-redis (v9)                                    │    │
│  │ • lib/pq (PostgreSQL)                                    │    │
│  │ • Dependency injection pattern                           │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                      External Services                           │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │ • Google Gemini (gemini-2.5-flash-preview)               │    │
│  │ • SerpAPI (Google Shopping)                              │    │
│  │ • Cloudflare Tunnel (SSL/CDN)                            │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                      Infrastructure                              │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │ • Docker 20.10+                                          │    │
│  │ • Docker Compose v2.0+                                   │    │
│  │ • PostgreSQL 18 Alpine                                   │    │
│  │ • Redis 8 Alpine                                         │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘
```

---

**Last Updated**: 2025-01-15
**Version**: 1.0.0
