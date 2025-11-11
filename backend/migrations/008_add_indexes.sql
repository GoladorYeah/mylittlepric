-- Migration: Add database indexes for performance optimization
-- Date: 2025-11-11
-- Purpose: Improve query performance for frequently accessed fields

-- User indexes
-- Index for email lookups during login (email is already unique, but this improves lookup performance)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index for OAuth provider lookups (Google login)
CREATE INDEX IF NOT EXISTS idx_users_provider_google_id ON users(provider, google_id);

-- SearchHistory indexes
-- Index for GetUserSearchHistory - filtering by user and ordering by date
CREATE INDEX IF NOT EXISTS idx_search_history_user_created ON search_history(user_id, created_at DESC);

-- Index for anonymous users - filtering by session and ordering by date
CREATE INDEX IF NOT EXISTS idx_search_history_session_created ON search_history(session_id, created_at DESC);

-- Partial index for cleanup job - only for anonymous users (user_id IS NULL)
-- This optimizes the CleanupExpiredAnonymousHistory query
CREATE INDEX IF NOT EXISTS idx_search_history_expires_anonymous
ON search_history(expires_at)
WHERE user_id IS NULL;

-- ChatSession indexes
-- Index for GetActiveSessionForUser - filtering by user and checking expiry
CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_expires ON chat_sessions(user_id, expires_at DESC);

-- Index for cleanup operations - finding expired sessions
CREATE INDEX IF NOT EXISTS idx_chat_sessions_expires ON chat_sessions(expires_at);

-- Index for session_id lookups (frequently used in ProcessChat)
-- Note: session_id is already unique, but adding index for performance
CREATE INDEX IF NOT EXISTS idx_chat_sessions_session_id ON chat_sessions(session_id);

-- Verify indexes were created
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
    AND tablename IN ('users', 'search_history', 'chat_sessions')
ORDER BY tablename, indexname;
