# Database Migrations

This directory contains SQL migration files for the MyLittlePrice database schema.

## Migration Files

Migrations are numbered sequentially and should be applied in order:

1. `001_initial_schema.sql` - Initial database schema with chat_sessions table
2. `002_add_users.sql` - User authentication tables
3. `003_add_search_history.sql` - Search history tracking
4. `004_fix_search_history_constraints.sql` - Constraint fixes
5. **`005_change_session_id_to_text.sql`** - **FIX: Changes session_id from UUID to TEXT**
6. `006_add_oauth_fields.sql` - OAuth integration fields
7. `007_add_session_state_fields.sql` - Session state tracking

## The UUID Validation Error Fix

### Problem
The application was throwing this error:
```
ERROR: invalid input syntax for type uuid: "1762698308128-tgjil1ulz"
```

This occurred because:
- The `search_history.session_id` column was defined as `UUID` type in migration 003
- The application generates custom session IDs like `"1762698308128-tgjil1ulz"` (timestamp-based)
- These custom IDs are not valid UUID format

### Solution
Migration `005_change_session_id_to_text.sql` changes the column type from UUID to TEXT to support the custom session ID format.

## How to Apply Migrations

### Option 1: Using Docker (Recommended)

From the project root directory:

```bash
# Start the database if not running
./docker.sh up

# Apply all pending migrations
./docker.sh db-migrate
```

### Option 2: Using the Migration Script Directly

If you have direct access to the database:

```bash
cd backend/migrations

# Set environment variables (if needed)
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=mylittleprice
export DB_USER=postgres
export DB_PASSWORD=postgres

# Run the migration script
./apply_migrations.sh
```

### Option 3: Manual Database Shell

```bash
# Open database shell
./docker.sh db-shell

# Or if using psql directly:
# psql -h localhost -U postgres -d mylittleprice

# Then run the migration manually:
\i /docker-entrypoint-initdb.d/005_change_session_id_to_text.sql
```

### Option 4: From Docker Container

```bash
# Execute directly in the postgres container
docker-compose exec postgres psql -U postgres -d mylittleprice \
  -f /docker-entrypoint-initdb.d/005_change_session_id_to_text.sql
```

## Migration Tracking

The `apply_migrations.sh` script creates a `schema_migrations` table to track which migrations have been applied:

```sql
CREATE TABLE schema_migrations (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

This ensures migrations are only applied once and can be safely re-run.

## Verifying the Fix

After applying migration 005, verify the change:

```sql
-- Check the column type
\d search_history

-- Should show:
-- session_id | text |
```

You can also test by creating a search history entry with a custom session ID:

```sql
INSERT INTO search_history (
    session_id, search_query, search_type,
    country_code, language_code, currency
) VALUES (
    '1762698308128-tgjil1ulz',
    'test query',
    'category',
    'US', 'en', 'USD'
);
```

This should now work without the UUID validation error.

## Important Notes

### Docker Volume Initialization

The postgres container mounts this directory at `/docker-entrypoint-initdb.d/`:

```yaml
volumes:
  - ./backend/migrations:/docker-entrypoint-initdb.d
```

**⚠️ Important**: Files in `/docker-entrypoint-initdb.d/` only run automatically when the database is **first created** (when the volume is empty).

If you added new migrations after the container was already created, you must:
1. Run migrations manually using one of the methods above, OR
2. Destroy and recreate the database volume:
   ```bash
   docker-compose down -v  # ⚠️ WARNING: This deletes all data!
   docker-compose up
   ```

### Production Deployment

For production environments:

1. **Always backup** the database before running migrations
2. Test migrations in a staging environment first
3. Run migrations during a maintenance window if possible
4. Check application logs after migration to ensure no errors

### Adding New Migrations

When creating a new migration:

1. Name it with the next sequential number: `00X_description.sql`
2. Add it to this README with a description
3. Test it locally first
4. Include both UP and DOWN operations if possible (for rollback)
5. Document any breaking changes or required application updates

## Troubleshooting

### "Permission denied" when running apply_migrations.sh

```bash
chmod +x backend/migrations/apply_migrations.sh
```

### "connection refused" when running migrations

Make sure the database container is running:
```bash
docker-compose ps
```

### Migration already applied but showing error

Check the schema_migrations table:
```sql
SELECT * FROM schema_migrations ORDER BY applied_at DESC;
```

### Need to revert a migration

Currently there are no automated rollback scripts. To revert:
1. Backup your database
2. Write and test a revert SQL script
3. Apply manually with caution
