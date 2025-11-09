-- migrations/009_add_search_sync.sql
-- Add search synchronization support for cross-device continuity

-- Add last_active_session_id to user_preferences for tracking ongoing searches
ALTER TABLE user_preferences
ADD COLUMN IF NOT EXISTS last_active_session_id TEXT;

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_user_preferences_last_active_session
ON user_preferences(last_active_session_id);

-- Add comment explaining the feature
COMMENT ON COLUMN user_preferences.last_active_session_id IS
'Tracks the most recent session with an unfinished search. Used for cross-device search continuity.';
