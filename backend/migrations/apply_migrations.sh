#!/bin/bash

# Migration Runner Script
# This script applies all pending SQL migrations to the database

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database connection parameters (can be overridden by environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-mylittleprice}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

MIGRATIONS_DIR="$(cd "$(dirname "$0")" && pwd)"

echo -e "${BLUE}🔍 Database Migration Runner${NC}"
echo "================================"
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo "Migrations dir: $MIGRATIONS_DIR"
echo ""

# Test database connection
echo -e "${YELLOW}Testing database connection...${NC}"
export PGPASSWORD="$DB_PASSWORD"
if ! psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}❌ Failed to connect to database${NC}"
    echo "Please check your connection parameters."
    exit 1
fi
echo -e "${GREEN}✅ Connected successfully${NC}"
echo ""

# Create migration tracking table if it doesn't exist
echo -e "${YELLOW}Creating migration tracking table...${NC}"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF
echo -e "${GREEN}✅ Migration tracking table ready${NC}"
echo ""

# Get list of already applied migrations
APPLIED_MIGRATIONS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT filename FROM schema_migrations ORDER BY filename;")

# Find all migration files
MIGRATION_FILES=$(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort)

if [ -z "$MIGRATION_FILES" ]; then
    echo -e "${YELLOW}⚠️  No migration files found in $MIGRATIONS_DIR${NC}"
    exit 0
fi

# Apply each migration
APPLIED_COUNT=0
SKIPPED_COUNT=0

for migration_file in $MIGRATION_FILES; do
    filename=$(basename "$migration_file")

    # Check if migration has already been applied
    if echo "$APPLIED_MIGRATIONS" | grep -q "$filename"; then
        echo -e "${BLUE}⏭️  Skipping $filename (already applied)${NC}"
        SKIPPED_COUNT=$((SKIPPED_COUNT + 1))
        continue
    fi

    echo -e "${YELLOW}📝 Applying $filename...${NC}"

    # Apply the migration
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"; then
        # Record successful migration
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c \
            "INSERT INTO schema_migrations (filename) VALUES ('$filename');"
        echo -e "${GREEN}✅ Applied $filename${NC}"
        APPLIED_COUNT=$((APPLIED_COUNT + 1))
    else
        echo -e "${RED}❌ Failed to apply $filename${NC}"
        exit 1
    fi
    echo ""
done

# Summary
echo "================================"
echo -e "${GREEN}✅ Migration complete!${NC}"
echo "Applied: $APPLIED_COUNT"
echo "Skipped: $SKIPPED_COUNT"
echo ""

unset PGPASSWORD
