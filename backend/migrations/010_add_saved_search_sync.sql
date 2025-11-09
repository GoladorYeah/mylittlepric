-- migrations/010_add_saved_search_sync.sql
-- Add saved_search field for cross-device "Last search saved" synchronization

-- Add saved_search column to user_preferences
ALTER TABLE user_preferences
ADD COLUMN IF NOT EXISTS saved_search JSONB;

-- Add index for faster JSON queries (optional, but helpful for large datasets)
CREATE INDEX IF NOT EXISTS idx_user_preferences_saved_search
ON user_preferences USING GIN (saved_search);

-- Add comment explaining the feature
COMMENT ON COLUMN user_preferences.saved_search IS
'Stores the last saved search state (before clicking "New Search"). Syncs across devices to show "Last search saved" banner.';
