-- Add state fields to chat_sessions table for full session persistence

-- Add search_state as JSONB
ALTER TABLE chat_sessions ADD COLUMN IF NOT EXISTS search_state JSONB DEFAULT '{"status":"idle","category":"","search_count":0}'::jsonb;

-- Add cycle_state as JSONB
ALTER TABLE chat_sessions ADD COLUMN IF NOT EXISTS cycle_state JSONB DEFAULT '{"cycle_id":1,"iteration":1,"cycle_history":[],"last_defined":[],"prompt_id":"UniversalPrompt v1.0.1","prompt_hash":""}'::jsonb;

-- Add conversation_context as JSONB (optional, can be NULL)
ALTER TABLE chat_sessions ADD COLUMN IF NOT EXISTS conversation_context JSONB;

-- Add user_id to allow linking sessions to authenticated users (optional)
ALTER TABLE chat_sessions ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- Update messages table to store products
ALTER TABLE messages ADD COLUMN IF NOT EXISTS products JSONB;
ALTER TABLE messages ADD COLUMN IF NOT EXISTS search_info JSONB;

-- Create index for user_id lookups
CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_id ON chat_sessions(user_id);

-- Create index for session_id lookups (for quick access)
CREATE INDEX IF NOT EXISTS idx_chat_sessions_session_id_text ON chat_sessions(session_id);

-- Add comment explaining the schema
COMMENT ON COLUMN chat_sessions.search_state IS 'Stores current search state (status, category, last_product, search_count)';
COMMENT ON COLUMN chat_sessions.cycle_state IS 'Stores Universal Prompt cycle state (cycle_id, iteration, history, context)';
COMMENT ON COLUMN chat_sessions.conversation_context IS 'Stores optimized conversation context for token efficiency (summary, preferences, exclusions)';
COMMENT ON COLUMN chat_sessions.user_id IS 'Optional link to authenticated user for session claiming';
