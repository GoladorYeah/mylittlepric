-- migrations/003_add_search_history.sql

-- Search history table for authenticated and anonymous users
CREATE TABLE IF NOT EXISTS search_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,  -- NULL for anonymous users
    session_id UUID REFERENCES chat_sessions(id) ON DELETE SET NULL,

    -- Search details
    search_query TEXT NOT NULL,
    optimized_query TEXT,                      -- Translated/optimized query used for API
    search_type VARCHAR(20) NOT NULL,          -- 'exact', 'parameters', 'category'
    category VARCHAR(100),                     -- Product category

    -- Search parameters
    country_code VARCHAR(2) NOT NULL,
    language_code VARCHAR(5) NOT NULL,
    currency VARCHAR(3) NOT NULL,

    -- Results metadata
    result_count INT DEFAULT 0,
    products_found JSONB,                      -- Array of product cards (limited to top 3-5)

    -- Interaction tracking
    clicked_product_id TEXT,                   -- Page token of clicked product
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- For cleanup of anonymous searches
    expires_at TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_search_history_user_id ON search_history(user_id);
CREATE INDEX idx_search_history_session_id ON search_history(session_id);
CREATE INDEX idx_search_history_created_at ON search_history(created_at DESC);
CREATE INDEX idx_search_history_expires_at ON search_history(expires_at) WHERE expires_at IS NOT NULL;

-- Add user_id to chat_sessions for authenticated users
ALTER TABLE chat_sessions ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_id ON chat_sessions(user_id);

-- Cleanup old anonymous search history (will be handled by app or cron)
-- DELETE FROM search_history WHERE expires_at < NOW() AND user_id IS NULL;
