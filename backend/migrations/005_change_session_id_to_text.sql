-- Change session_id from UUID to TEXT to support custom session ID format

-- Drop foreign key constraint first
ALTER TABLE search_history DROP CONSTRAINT IF EXISTS search_history_session_id_fkey;

-- Change column type
ALTER TABLE search_history ALTER COLUMN session_id TYPE TEXT USING session_id::TEXT;

-- Recreate index
DROP INDEX IF EXISTS idx_search_history_session_id;
CREATE INDEX idx_search_history_session_id ON search_history(session_id);
