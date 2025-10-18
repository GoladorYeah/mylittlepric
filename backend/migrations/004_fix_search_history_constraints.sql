-- migrations/004_fix_search_history_constraints.sql

-- Remove foreign key constraints that conflict with Redis-based sessions
-- Sessions are stored in Redis, not PostgreSQL, so we can't have FK constraint

-- Drop the session_id foreign key constraint
ALTER TABLE search_history
DROP CONSTRAINT IF EXISTS search_history_session_id_fkey;

-- Drop the user_id foreign key constraint and recreate without it
ALTER TABLE search_history
DROP CONSTRAINT IF EXISTS search_history_user_id_fkey;

-- Note: user_id and session_id remain as UUID columns but without FK constraints
-- This allows:
-- 1. Anonymous users: user_id = NULL, session_id = <redis session uuid>
-- 2. Authenticated users: user_id = <user uuid>, session_id = <redis session uuid>
-- 3. No database integrity checks - application logic handles validation

-- Add comment to clarify the design
COMMENT ON COLUMN search_history.user_id IS 'UUID of authenticated user (from users table), NULL for anonymous users. No FK constraint - validated by application.';
COMMENT ON COLUMN search_history.session_id IS 'UUID of chat session (stored in Redis), not in PostgreSQL. No FK constraint - sessions are ephemeral.';
