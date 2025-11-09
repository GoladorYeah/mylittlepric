#!/bin/bash

# Quick migration script to apply pending database migrations
# This script works with the existing Docker setup

set -e

echo "🔍 Applying database migrations..."
echo ""

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: docker-compose is not available"
    echo "Please install Docker and docker-compose first"
    exit 1
fi

# Check if containers are running
if ! docker-compose ps | grep -q "mylittleprice-postgres.*Up"; then
    echo "❌ Error: PostgreSQL container is not running"
    echo "Please start the services first with: ./docker.sh up"
    exit 1
fi

# Run migrations using the docker.sh script
echo "📊 Running migrations in PostgreSQL container..."
docker-compose exec -T postgres /bin/sh -c "
    export PGPASSWORD=postgres
    cd /docker-entrypoint-initdb.d

    # Create migrations tracking table
    psql -U postgres -d mylittleprice -c \"
        CREATE TABLE IF NOT EXISTS schema_migrations (
            id SERIAL PRIMARY KEY,
            filename VARCHAR(255) NOT NULL UNIQUE,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    \"

    # Get applied migrations
    APPLIED=\$(psql -U postgres -d mylittleprice -t -c 'SELECT filename FROM schema_migrations ORDER BY filename;')

    # Apply each migration file in order
    for f in \$(ls -1 *.sql 2>/dev/null | sort); do
        if echo \"\$APPLIED\" | grep -q \"\$f\"; then
            echo \"⏭️  Skipping \$f (already applied)\"
        else
            echo \"📝 Applying \$f...\"
            if psql -U postgres -d mylittleprice -f \"\$f\"; then
                psql -U postgres -d mylittleprice -c \"INSERT INTO schema_migrations (filename) VALUES ('\$f');\"
                echo \"✅ Applied \$f\"
            else
                echo \"❌ Failed to apply \$f\"
                exit 1
            fi
        fi
    done
"

echo ""
echo "✅ Migration complete!"
echo ""
echo "Verifying picture column exists..."
docker-compose exec -T postgres psql -U postgres -d mylittleprice -c "\d users" | grep picture && echo "✅ Picture column found!" || echo "❌ Picture column not found"
