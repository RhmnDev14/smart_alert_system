-- User memories table for persistent memory (OpenClaw-style)
-- Stores facts, preferences, and habits that the AI learns about each user
CREATE TABLE IF NOT EXISTS user_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    memory_type VARCHAR(50) NOT NULL DEFAULT 'fact', -- 'preference', 'fact', 'habit', 'personal'
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for fast lookup by user
CREATE INDEX IF NOT EXISTS idx_user_memories_user_id ON user_memories(user_id);
