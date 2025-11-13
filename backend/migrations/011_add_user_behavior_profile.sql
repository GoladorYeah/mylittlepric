-- migrations/011_add_user_behavior_profile.sql
-- User behavior profile for long-term preference tracking and personalized recommendations

-- User behavior profile table
CREATE TABLE IF NOT EXISTS user_behavior_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Category preferences with weights (JSONB: {"electronics": 0.8, "clothing": 0.6})
    category_preferences JSONB DEFAULT '{}'::jsonb,

    -- Price ranges by category (JSONB: {"electronics": {"min": 100, "max": 1000}})
    price_ranges JSONB DEFAULT '{}'::jsonb,

    -- Preferred brands with frequency counts (JSONB: {"Apple": 15, "Samsung": 8})
    brand_preferences JSONB DEFAULT '{}'::jsonb,

    -- Communication style: "brief", "balanced", "detailed"
    communication_style VARCHAR(20) DEFAULT 'balanced',

    -- Preferred search type: "product", "specification", "comparison"
    preferred_search_type VARCHAR(50),

    -- Session metrics
    avg_session_duration DOUBLE PRECISION DEFAULT 0,
    avg_messages_per_session DOUBLE PRECISION DEFAULT 0,
    success_rate DOUBLE PRECISION DEFAULT 0,

    -- Common keywords extracted from user messages
    common_keywords JSONB DEFAULT '[]'::jsonb,

    -- Time-based patterns (preferred shopping hours)
    time_patterns JSONB DEFAULT '{}'::jsonb,

    -- Counters
    total_sessions INTEGER DEFAULT 0,
    total_products_viewed INTEGER DEFAULT 0,
    total_products_clicked INTEGER DEFAULT 0,

    -- Timestamps
    last_learning_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_behavior_profiles_user_id
    ON user_behavior_profiles(user_id);

CREATE INDEX IF NOT EXISTS idx_user_behavior_profiles_learning_update
    ON user_behavior_profiles(last_learning_update);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_user_behavior_profiles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at
CREATE TRIGGER trigger_user_behavior_profiles_updated_at
    BEFORE UPDATE ON user_behavior_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_user_behavior_profiles_updated_at();

-- Create default behavior profile for existing users
INSERT INTO user_behavior_profiles (user_id)
SELECT id FROM users
ON CONFLICT (user_id) DO NOTHING;
