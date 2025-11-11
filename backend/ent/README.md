# Ent ORM Schema

This directory contains the Ent ORM schema definitions and generated code for MyLittlePrice.

## 📁 Structure

```
ent/
├── schema/              # Schema definitions (edit these)
│   ├── user.go
│   ├── chatsession.go
│   ├── message.go
│   ├── searchhistory.go
│   └── userpreference.go
├── generate.go          # Code generation trigger
└── [generated files]    # Auto-generated code (DO NOT EDIT)
```

## 🔨 Commands

### Generate Code
After modifying schema files:
```bash
cd backend
go generate ./ent
```

### Run Migrations
Automatically create/update database tables:
```go
client.Schema.Create(ctx)
```

### Reset Database (Development Only)
```bash
docker exec mylittleprice-postgres psql -U postgres -d mylittleprice -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
```

## 🗂️ Entities

### User
- Email/Google OAuth authentication
- Relationships: sessions, search_history, preferences

### ChatSession
- Session management with UUID
- JSONB fields: search_state, cycle_state, conversation_context
- Relationships: user, messages

### Message
- Chat messages (user/assistant)
- JSONB fields: products, search_info
- String array: quick_replies

### SearchHistory
- User search tracking
- JSONB: products_found

### UserPreference
- User settings (country, language, currency, theme)
- One-to-one with User

## 📖 Usage Examples

### Create User
```go
user, err := client.User.
    Create().
    SetEmail("user@example.com").
    SetName("John Doe").
    SetProvider("email").
    Save(ctx)
```

### Query with Relations
```go
user, err := client.User.
    Query().
    Where(user.Email("user@example.com")).
    WithSessions().  // Eager load sessions
    WithPreferences(). // Eager load preferences
    Only(ctx)
```

### Create Session with User
```go
session, err := client.ChatSession.
    Create().
    SetSessionID("session-123").
    SetUser(user).
    SetCountryCode("US").
    SetSearchState(map[string]interface{}{
        "status": "idle",
    }).
    Save(ctx)
```

### Update
```go
user, err := client.User.
    UpdateOneID(userID).
    SetName("New Name").
    Save(ctx)
```

### Delete
```go
err := client.User.
    DeleteOneID(userID).
    Exec(ctx)
```

## 🔗 Relationships

```
User 1─────∞ ChatSession
User 1─────∞ SearchHistory
User 1─────1 UserPreference

ChatSession 1─────∞ Message
```

## 📚 Resources

- [Ent Documentation](https://entgo.io/)
- [Schema Fields](https://entgo.io/docs/schema-fields)
- [Schema Edges](https://entgo.io/docs/schema-edges)
- [CRUD Operations](https://entgo.io/docs/crud)
