# Environment Configuration Guide

## 📁 File Structure

```
mylittleprice/
├── .env.example          ← Main template (USE THIS!)
├── .env                  ← Your actual config (create from .env.example)
│
├── backend/
│   ├── .env.example      ← Only for local Go development without Docker
│   └── .env              ← Only used when running "go run" locally
│
└── frontend/
    └── .env              ← Only for local Next.js development without Docker
```

## 🎯 Which file to use?

### Using Docker (Recommended)
**You only need ONE file:**

1. Copy the root `.env.example` to `.env`
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` and fill in:
   - `GEMINI_API_KEYS` - Get from https://makersuite.google.com/app/apikey
   - `SERP_API_KEYS` - Get from https://serpapi.com/
   - `JWT_ACCESS_SECRET` - Generate with `openssl rand -hex 32`
   - `JWT_REFRESH_SECRET` - Generate with `openssl rand -hex 32`

3. Run Docker:
   ```bash
   ./docker.sh up
   ```

**That's it!** The root `.env` file is used by both development and production Docker setups.

### Local Development (Without Docker)

If you're running backend and frontend separately on your machine:

1. **Backend**: Copy `backend/.env.example` to `backend/.env`
2. **Frontend**: Create `frontend/.env` with:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

## 🔐 Security Checklist

Before deploying to production, make sure you've changed:

- [ ] `JWT_ACCESS_SECRET` - Must be unique and random
- [ ] `JWT_REFRESH_SECRET` - Must be different from access secret
- [ ] `POSTGRES_PASSWORD` - Strong password
- [ ] `REDIS_PASSWORD` - Strong password (optional but recommended)
- [ ] `CORS_ORIGINS` - Set to your actual domain
- [ ] `GEMINI_API_KEYS` - Real API keys
- [ ] `SERP_API_KEYS` - Real API keys

### Generate Secure Secrets

```bash
# Generate JWT secrets (run twice for access & refresh)
openssl rand -hex 32

# Or use this one-liner to update .env automatically
echo "JWT_ACCESS_SECRET=$(openssl rand -hex 32)" >> .env
echo "JWT_REFRESH_SECRET=$(openssl rand -hex 32)" >> .env
```

## 🚀 Quick Start Examples

### Development
```bash
# 1. Copy template
cp .env.example .env

# 2. Add your API keys to .env
nano .env  # or use any editor

# 3. Start everything
./docker.sh up

# Done! Frontend: http://localhost:3000
```

### Production
```bash
# 1. Copy template
cp .env.example .env

# 2. Edit .env and set:
#    - Real domain names (CORS_ORIGINS)
#    - Strong passwords
#    - Secure JWT secrets
#    - Real API keys

# 3. Set ENV=production
echo "ENV=production" >> .env

# 4. Deploy
./docker.sh prod-up
```

## ❌ Common Mistakes

### ❌ DON'T: Create multiple .env files
```
.env
.env.development
.env.production  ← Confusing!
.env.local
```

### ✅ DO: Use one .env file
```
.env  ← Single source of truth
```

### ❌ DON'T: Commit .env files
```bash
git add .env  # NEVER DO THIS!
```

### ✅ DO: Use .env.example as template
```bash
git add .env.example  # This is OK - it's a template
```

## 📝 Environment Variables Reference

See [.env.example](.env.example) for the complete list of available variables with descriptions.

### Required Variables (Docker)
- `GEMINI_API_KEYS` - Gemini AI API keys
- `SERP_API_KEYS` - SERP API keys for product search

### Required for Production
- `JWT_ACCESS_SECRET` - JWT signing secret
- `JWT_REFRESH_SECRET` - JWT refresh token secret
- `POSTGRES_PASSWORD` - Database password
- `CORS_ORIGINS` - Allowed domains

### Optional (have defaults)
- `BACKEND_PORT` - Default: 8080
- `FRONTEND_PORT` - Default: 3000
- `SESSION_TTL` - Default: 86400 (24 hours)
- `MAX_MESSAGES_PER_SESSION` - Default: 8
- And many more...

## 🆘 Troubleshooting

### "GEMINI_API_KEYS is required"
→ Make sure you've added your Gemini API key to `.env`

### "JWT_ACCESS_SECRET is required"
→ Generate and add JWT secrets to `.env`:
```bash
openssl rand -hex 32
```

### Docker not reading .env
→ Make sure the file is named exactly `.env` (not `.env.txt`)
→ Check file is in the root directory (same level as docker-compose.yml)

### Variables not working in production
→ Make sure you're using `./docker.sh prod-up` (not `docker-compose up`)
→ This command combines both docker-compose.yml and docker-compose.prod.yml

## 🎓 Understanding the System

```
┌─────────────────────────────────────────────┐
│  .env (root)                                │
│  Single source of truth for all env vars   │
└─────────────────┬───────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌───────────────┐   ┌──────────────────┐
│ docker-       │   │ docker-compose.  │
│ compose.yml   │   │ prod.yml         │
│               │   │                  │
│ Development   │   │ Production       │
│ + Defaults    │   │ + Overrides      │
└───────────────┘   └──────────────────┘
```

The production setup **combines** both files:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
```

This is why `./docker.sh prod-up` works perfectly!
