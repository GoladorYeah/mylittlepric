# MyLittlePrice 🛒

AI-powered product search assistant with conversational chat interface. Built with Go, Next.js, Google Gemini AI, and SerpAPI.

## ✨ Features

- 🤖 **Conversational AI** - Natural language product search powered by Google Gemini
- 🔍 **Smart Search** - Intelligent product matching with relevance scoring
- 🔐 **User Authentication** - ChatGPT-style auth with anonymous sessions
- 💬 **Real-time Chat** - WebSocket-based instant messaging
- 🌍 **Multi-language** - Automatic query translation and locale detection
- 📱 **Responsive UI** - Dark mode support with modern design
- 🐳 **Docker Ready** - Full containerization with Docker Compose

## 🚀 Quick Start (Docker)

1. **Clone and configure**
   ```bash
   git clone <your-repo>
   cd mylittleprice
   cp .env.example .env
   ```

2. **Add API keys to `.env`**
   - Get Gemini API key: https://makersuite.google.com/app/apikey
   - Get SERP API key: https://serpapi.com/
   - Add them to `.env` file

3. **Start everything**
   ```bash
   # Linux/Mac
   ./docker.sh up

   # Windows
   .\docker.ps1 up
   ```

4. **Open in browser**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Docs: http://localhost:8080/health

That's it! 🎉

## 📚 Documentation

- **[Environment Setup](ENV_SETUP.md)** - Detailed .env configuration guide
- **[Architecture](CLAUDE.md)** - Project architecture and technical details
- **[Docker Guide](DOCKER.md)** - Docker commands and deployment

## 🛠️ Tech Stack

### Backend
- **Go 1.23** - Fast, reliable API server
- **Fiber v2** - Express-like web framework
- **Google Gemini** - AI conversation & search
- **SerpAPI** - Google Shopping product data
- **Redis** - Caching and session storage
- **PostgreSQL** - User data and persistence
- **JWT** - Secure authentication

### Frontend
- **Next.js 15** - React framework with App Router
- **Bun** - Fast JavaScript runtime
- **Zustand** - State management
- **Tailwind CSS v4** - Utility-first styling
- **WebSocket** - Real-time communication
- **shadcn/ui** - UI components

## 🐳 Docker Commands

```bash
# Development
./docker.sh up          # Start all services
./docker.sh logs        # View logs
./docker.sh down        # Stop all services
./docker.sh build       # Rebuild containers

# Production
./docker.sh prod-up     # Start in production mode
./docker.sh prod-down   # Stop production
./docker.sh prod-logs   # View production logs

# Database
./docker.sh db-migrate  # Run migrations
./docker.sh db-shell    # Open PostgreSQL shell
./docker.sh redis-cli   # Open Redis CLI
```

## 🔧 Local Development (Without Docker)

### Backend
```bash
cd backend
cp .env.example .env
# Add your API keys to .env
go run cmd/api/main.go
```

### Frontend
```bash
cd frontend
bun install
bun run dev
```

## 📦 Project Structure

```
mylittleprice/
├── backend/              # Go API server
│   ├── cmd/api/         # Main application
│   ├── internal/        # Business logic
│   │   ├── handlers/   # HTTP & WebSocket handlers
│   │   ├── services/   # Core services (Gemini, SERP, Auth)
│   │   ├── models/     # Data models
│   │   └── middleware/ # Auth middleware
│   └── migrations/     # Database schemas
│
├── frontend/            # Next.js application
│   ├── src/
│   │   ├── app/        # Pages (App Router)
│   │   ├── components/ # React components
│   │   └── lib/        # Utilities & stores
│   └── public/         # Static assets
│
├── docker-compose.yml       # Development setup
├── docker-compose.prod.yml  # Production overrides
├── .env.example            # Environment template
└── ENV_SETUP.md           # Configuration guide
```

## 🔐 Authentication

MyLittlePrice uses ChatGPT-style authentication:

- ✅ **Anonymous usage** - Start chatting immediately
- ✅ **Optional signup** - Create account to save history
- ✅ **Session migration** - Anonymous sessions automatically linked when you sign up
- ✅ **JWT tokens** - Secure access & refresh token flow

API Endpoints:
- `POST /api/auth/signup` - Create account
- `POST /api/auth/login` - Sign in
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Sign out
- `GET /api/auth/me` - Get user info
- `POST /api/auth/claim-sessions` - Link anonymous sessions

## 🌍 Environment Variables

See [ENV_SETUP.md](ENV_SETUP.md) for detailed configuration guide.

**Required:**
- `GEMINI_API_KEYS` - Gemini AI API keys
- `SERP_API_KEYS` - SERP API keys

**Production:**
- `JWT_ACCESS_SECRET` - Generate with `openssl rand -hex 32`
- `JWT_REFRESH_SECRET` - Different from access secret
- `POSTGRES_PASSWORD` - Strong password
- `CORS_ORIGINS` - Your domain(s)

## 📊 API Endpoints

### Chat
- `POST /api/chat` - Send message (HTTP)
- `GET /ws` - WebSocket chat endpoint

### Products
- `POST /api/product-details` - Get product details by token

### Statistics
- `GET /api/stats/keys` - API key usage
- `GET /api/stats/grounding` - Grounding decisions
- `GET /api/stats/tokens` - Token usage
- `GET /api/stats/all` - All statistics

### Health
- `GET /health` - Service health check

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 🆘 Support

- **Issues**: Open an issue on GitHub
- **Documentation**: Check [ENV_SETUP.md](ENV_SETUP.md) and [CLAUDE.md](CLAUDE.md)
- **API Keys**:
  - Gemini: https://makersuite.google.com/app/apikey
  - SERP: https://serpapi.com/

---

Made with ❤️ using Claude Code
