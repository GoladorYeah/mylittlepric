# Fix: Missing Picture Column in Users Table

## Problem

The application is failing with the following error when trying to sync users from Redis to PostgreSQL:

```
❌ Error upserting preferences for user: failed to sync user to PostgreSQL:
failed to save user to PostgreSQL: pq: column "picture" of relation "users" does not exist
```

## Root Cause

The database schema is missing the `picture` column in the `users` table. This column is required for OAuth authentication (Google, etc.) to store the user's profile picture URL.

The migration file `backend/migrations/006_add_oauth_fields.sql` exists and contains the fix, but it hasn't been applied to the database yet.

## Solution

### Option 1: Using the Helper Script (Recommended)

Run the migration helper script from the project root:

```bash
./apply-migrations.sh
```

This script will:
- Check if Docker containers are running
- Apply all pending migrations in order
- Track applied migrations to prevent duplicates
- Verify the picture column was created successfully

### Option 2: Using the Docker Management Script

Use the built-in docker.sh script:

```bash
./docker.sh db-migrate
```

### Option 3: Manual Application

If you prefer to apply the specific migration manually:

```bash
# Access the PostgreSQL container
docker-compose exec postgres psql -U postgres -d mylittleprice

# Run the following SQL:
ALTER TABLE users
ADD COLUMN IF NOT EXISTS picture VARCHAR(500),
ADD COLUMN IF NOT EXISTS provider VARCHAR(50) NOT NULL DEFAULT 'email',
ADD COLUMN IF NOT EXISTS provider_id VARCHAR(255);

ALTER TABLE users
ALTER COLUMN password_hash DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider, provider_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_unique
ON users(provider, provider_id)
WHERE provider_id IS NOT NULL;
```

## Verification

After applying the migration, verify the column exists:

```bash
docker-compose exec postgres psql -U postgres -d mylittleprice -c "\d users"
```

You should see the `picture` column listed with type `character varying(500)`.

## What This Migration Adds

The `006_add_oauth_fields.sql` migration adds support for OAuth authentication:

- `picture`: Stores the user's profile picture URL (VARCHAR 500)
- `provider`: Authentication provider (email, google, etc.) - defaults to 'email'
- `provider_id`: Unique ID from the OAuth provider
- Makes `password_hash` optional (for OAuth users who don't have passwords)
- Adds indexes for efficient OAuth lookups

## Prevention

To prevent this issue in the future:

1. Always run `./docker.sh db-migrate` after pulling new code that includes migrations
2. Check the `backend/migrations/` directory for new migration files
3. The migration tracking table `schema_migrations` keeps track of applied migrations

## Related Files

- Migration file: `backend/migrations/006_add_oauth_fields.sql`
- Migration script: `backend/migrations/apply_migrations.sh`
- Docker helper: `docker.sh` (db-migrate command)
- New helper script: `apply-migrations.sh` (project root)
