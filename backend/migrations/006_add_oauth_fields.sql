-- migrations/006_add_oauth_fields.sql
-- Add OAuth provider fields to users table

-- Add OAuth-related columns
ALTER TABLE users
ADD COLUMN IF NOT EXISTS picture VARCHAR(500),
ADD COLUMN IF NOT EXISTS provider VARCHAR(50) NOT NULL DEFAULT 'email',
ADD COLUMN IF NOT EXISTS provider_id VARCHAR(255);

-- Make password_hash optional for OAuth users
ALTER TABLE users
ALTER COLUMN password_hash DROP NOT NULL;

-- Create index on provider + provider_id for faster OAuth lookups
CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider, provider_id);

-- Add unique constraint on provider + provider_id combination
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_unique
ON users(provider, provider_id)
WHERE provider_id IS NOT NULL;
