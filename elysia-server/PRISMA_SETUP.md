# Prisma ORM Setup for Elysia Server

This document explains the Prisma ORM integration with the Elysia server.

## Overview

The elysia-server now uses **Prisma ORM** with **Prismabox** for type-safe database operations and automatic TypeBox schema generation for Elysia validation.

### Why Prisma?

- ✅ **Type-safe**: Auto-generated TypeScript types
- ✅ **Better DX**: Intuitive query API instead of raw SQL
- ✅ **Fixed bugs**: The old `pg` implementation had template literal syntax issues
- ✅ **Migrations**: Built-in schema migration support
- ✅ **Prismabox integration**: Auto-generates TypeBox schemas for Elysia validation

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Elysia Server                      │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────┐         ┌──────────────┐        │
│  │  Container   │────────▶│   Services   │        │
│  │   (DI)       │         │              │        │
│  └──────┬───────┘         └──────┬───────┘        │
│         │                        │                 │
│         │                        │                 │
│  ┌──────▼────────────────────────▼───┐            │
│  │      PrismaClient (Singleton)     │            │
│  └──────────────┬────────────────────┘            │
│                 │                                  │
└─────────────────┼──────────────────────────────────┘
                  │
                  ▼
         ┌────────────────┐
         │  PostgreSQL DB │
         └────────────────┘
```

## Database Schema

The Prisma schema defines 8 tables:

### Core Tables
1. **users** - User authentication with Google OAuth support
2. **refresh_tokens** - JWT token rotation
3. **chat_sessions** - Chat session state with JSONB columns
4. **messages** - Conversation history

### Analytics Tables
5. **search_history** - User search analytics
6. **search_queries** - Query optimization tracking
7. **products_cache** - Popular product caching
8. **api_usage** - API key usage statistics

## Setup Instructions

### 1. Install Dependencies

Dependencies are already installed:
- `@prisma/client` (6.19.0)
- `prisma` (6.19.0)
- `prismabox` (1.1.25)

### 2. Configure DATABASE_URL

Create a `.env` file (or add to existing):

```bash
DATABASE_URL="postgresql://user:password@localhost:5432/mylittleprice?schema=public"
```

### 3. Generate Prisma Client

```bash
npm run prisma:generate
```

This runs:
1. `prisma generate` - Generates Prisma Client
2. `prismabox` - Generates TypeBox schemas in `src/generated/prismabox/`

### 4. Sync Database Schema

**Option A: Push schema without migrations (development)**
```bash
npm run db:push
```

**Option B: Create migration (production)**
```bash
npm run prisma:migrate
```

### 5. Run the Server

```bash
npm run dev
```

## Usage Examples

### Accessing Prisma Client

Prisma client is available through the DI container:

```typescript
// In a service
export class MyService {
  constructor(private prisma: PrismaClient) {}

  async getUsers() {
    return await this.prisma.user.findMany();
  }
}
```

### SearchHistoryService Examples

```typescript
// Create search entry
await searchHistoryService.createEntry(
  userId,
  sessionId,
  query,
  category,
  searchType,
  country,
  language,
  resultCount
);

// Get user history
const history = await searchHistoryService.getUserHistory(userId, 20);

// Get session history
const sessionHistory = await searchHistoryService.getSessionHistory(sessionId);

// Get search statistics
const stats = await searchHistoryService.getUserSearchStats(userId);
```

### Direct Prisma Queries

```typescript
// Create a user
const user = await container.prisma.user.create({
  data: {
    email: 'user@example.com',
    passwordHash: hashedPassword,
    fullName: 'John Doe',
  },
});

// Find with relations
const session = await container.prisma.chatSession.findUnique({
  where: { sessionId: 'abc123' },
  include: {
    messages: true,
    user: true,
  },
});

// Update with JSONB
await container.prisma.chatSession.update({
  where: { id: sessionId },
  data: {
    searchState: {
      status: 'completed',
      category: 'electronics',
      search_count: 3,
    },
  },
});
```

## Prismabox Integration

Prismabox automatically generates TypeBox schemas from Prisma models for Elysia validation.

### Generated Schemas Location

```
src/generated/prismabox/
├── User.ts
├── ChatSession.ts
├── Message.ts
├── SearchHistory.ts
└── ... (all models)
```

### Using with Elysia Validation

```typescript
import { User } from './generated/prismabox/User';
import { Elysia, t } from 'elysia';

const app = new Elysia()
  .post('/users', async ({ body, set }) => {
    // body is automatically validated against User schema
    return await prisma.user.create({ data: body });
  }, {
    body: User,
    response: User
  });
```

## Available Scripts

| Command | Description |
|---------|-------------|
| `npm run prisma:generate` | Generate Prisma Client + Prismabox schemas |
| `npm run prisma:migrate` | Create and apply migration |
| `npm run prisma:studio` | Open Prisma Studio (GUI) |
| `npm run db:push` | Push schema changes without migration |
| `npm run dev` | Start development server |

## Prisma Studio

Prisma Studio provides a GUI for your database:

```bash
npm run prisma:studio
```

Opens at `http://localhost:5555`

## Migration from Old `pg` Implementation

### What Changed

**Before (Broken):**
```typescript
// ❌ This syntax doesn't work with pg.Pool
const [entry] = await this.sql`
  INSERT INTO search_history (...)
  VALUES (${id}, ${userId}, ...)
  RETURNING *
`;
```

**After (Working):**
```typescript
// ✅ Type-safe Prisma query
const entry = await this.prisma.searchHistory.create({
  data: {
    userId,
    sessionId,
    searchQuery: query,
    // ...
  },
});
```

### Breaking Changes

1. **SearchHistoryItem type**: Now uses Prisma's `SearchHistory` type
2. **Field names**: Camel case (e.g., `searchQuery` instead of `query`)
3. **Return types**: Prisma returns full objects with all fields

### Container Changes

```typescript
// Old
public db: any;
this.db = initPostgres(config);
await this.db.end();

// New
public prisma: PrismaClient;
this.prisma = initPrisma();
await this.prisma.$disconnect();
```

## JSONB Columns

Prisma handles JSONB columns seamlessly:

```typescript
// ChatSession has JSONB columns
const session = await prisma.chatSession.update({
  where: { sessionId: 'abc' },
  data: {
    searchState: {
      status: 'searching',
      category: 'laptops',
      search_count: 1,
    },
    cycleState: {
      cycle_id: 1,
      iteration: 2,
      cycle_history: [...],
    },
  },
});

// TypeScript knows the structure!
console.log(session.searchState.category); // ✅ Type-safe
```

## Performance Considerations

### Connection Pooling

Prisma manages connection pooling automatically. The singleton pattern prevents multiple instances:

```typescript
// Global singleton for hot reload
declare global {
  var prisma: PrismaClient | undefined;
}

export function initPrisma(): PrismaClient {
  if (global.prisma) return global.prisma;

  const prisma = new PrismaClient();
  if (process.env.NODE_ENV === 'development') {
    global.prisma = prisma;
  }
  return prisma;
}
```

### Query Optimization

```typescript
// Bad: N+1 query problem
const sessions = await prisma.chatSession.findMany();
for (const session of sessions) {
  const user = await prisma.user.findUnique({ where: { id: session.userId } });
}

// Good: Single query with include
const sessions = await prisma.chatSession.findMany({
  include: { user: true }
});
```

## Troubleshooting

### "PrismaClient is unable to be run in the browser"

- Make sure you're not importing Prisma Client in frontend code
- Prisma Client only works server-side

### "Cannot find module '@prisma/client'"

Run:
```bash
npm install
npm run prisma:generate
```

### Migration Issues

```bash
# Reset database (⚠️ destroys data)
npx prisma migrate reset

# Push without migration
npm run db:push
```

### Prisma Binary Download Errors

If you encounter 403 errors downloading Prisma binaries:
- Check network/firewall settings
- Try: `PRISMA_ENGINES_CHECKSUM_IGNORE_MISSING=1 npx prisma generate`
- Or use `prisma db push` which doesn't require migration files

## Resources

- [Prisma Documentation](https://www.prisma.io/docs)
- [Elysia Documentation](https://elysiajs.com)
- [Prismabox GitHub](https://github.com/m1212e/prismabox)
- [TypeBox Documentation](https://github.com/sinclairzx81/typebox)

## Next Steps

1. ✅ Prisma client configured
2. ✅ SearchHistoryService migrated
3. ✅ Container updated with DI
4. ⏳ Generate Prismabox schemas (when binaries available)
5. ⏳ Add validation with TypeBox schemas
6. ⏳ Implement remaining database services (SessionService, AuthService)

## Notes

- The old `pg` dependency has been removed from `package.json`
- All SQL template literal syntax has been replaced with Prisma queries
- Type safety is now enforced at compile time
- Database schema is version-controlled through Prisma migrations
