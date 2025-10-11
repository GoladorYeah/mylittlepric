-- migrations/001_initial_schema.sql

-- Chat sessions table
CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) UNIQUE NOT NULL,
    country_code VARCHAR(2) NOT NULL,        -- CH, DE, IT, etc.
    language_code VARCHAR(5) NOT NULL,        -- de, fr, it, en
    currency VARCHAR(3) NOT NULL,             -- CHF, EUR, etc.
    message_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,                -- 'user' or 'assistant'
    content TEXT NOT NULL,
    response_type VARCHAR(20),                -- 'dialogue', 'search', 'product_card'
    quick_replies JSONB,                      -- Array of quick reply options
    metadata JSONB,                           -- Additional data
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Search queries table (for analytics and caching)
CREATE TABLE IF NOT EXISTS search_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES chat_sessions(id) ON DELETE CASCADE,
    original_query TEXT NOT NULL,
    optimized_query TEXT NOT NULL,
    search_type VARCHAR(20) NOT NULL,         -- 'exact', 'parameters', 'category'
    country_code VARCHAR(2) NOT NULL,
    result_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products cache table (for popular products)
CREATE TABLE IF NOT EXISTS products_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_token TEXT UNIQUE NOT NULL,              -- Google Immersive Product token
    product_name TEXT NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    price_info JSONB NOT NULL,                    -- Price, currency, etc.
    merchant_info JSONB,                          -- Merchant data
    product_details JSONB,                        -- Full product data
    image_url TEXT,
    link TEXT,
    rating DECIMAL(3,2),
    fetch_count INT DEFAULT 0,                    -- Popularity metric
    last_fetched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- API usage tracking (for rotation and analytics)
CREATE TABLE IF NOT EXISTS api_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_name VARCHAR(50) NOT NULL,            -- 'gemini' or 'serp'
    key_index INT NOT NULL,                   -- Which key was used
    request_type VARCHAR(50),                 -- Type of request
    response_time_ms INT,
    success BOOLEAN DEFAULT true,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_chat_sessions_session_id ON chat_sessions(session_id);
CREATE INDEX idx_chat_sessions_expires_at ON chat_sessions(expires_at);
CREATE INDEX idx_messages_session_id ON messages(session_id);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX idx_search_queries_session_id ON search_queries(session_id);
CREATE INDEX idx_search_queries_optimized ON search_queries(optimized_query);
CREATE INDEX idx_products_cache_cid ON products_cache(product_cid);
CREATE INDEX idx_products_cache_popularity ON products_cache(fetch_count DESC);
CREATE INDEX idx_api_usage_created_at ON api_usage(created_at DESC);

-- Cleanup old sessions (will be handled by app or cron)
-- DELETE FROM chat_sessions WHERE expires_at < NOW();