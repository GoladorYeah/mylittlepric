-- migrations/008_add_user_preferences.sql
-- User preferences for cross-device sync (country, currency, language, theme, UI settings)

-- User preferences table
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Regional settings
    country VARCHAR(2),  -- ISO 3166-1 alpha-2 country code (e.g., "US", "GB")
    currency VARCHAR(3), -- ISO 4217 currency code (e.g., "USD", "EUR")
    language VARCHAR(5), -- ISO 639-1 language code (e.g., "en", "es")

    -- UI settings
    theme VARCHAR(20),   -- "light", "dark", "system"
    sidebar_open BOOLEAN DEFAULT true,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_user_preferences_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at
CREATE TRIGGER trigger_user_preferences_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_user_preferences_updated_at();
