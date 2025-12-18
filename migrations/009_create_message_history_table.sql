-- Create MESSAGE_HISTORY table
CREATE TABLE IF NOT EXISTS message_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_content TEXT NOT NULL,
    message_type VARCHAR(20) NOT NULL CHECK (message_type IN ('incoming', 'outgoing')),
    intent_detected VARCHAR(100),
    ai_response TEXT,
    received_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    is_processed BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_message_history_user_id ON message_history(user_id);
CREATE INDEX IF NOT EXISTS idx_message_history_received_at ON message_history(received_at);
CREATE INDEX IF NOT EXISTS idx_message_history_message_type ON message_history(message_type);
CREATE INDEX IF NOT EXISTS idx_message_history_is_processed ON message_history(is_processed);
CREATE INDEX IF NOT EXISTS idx_message_history_user_received_at ON message_history(user_id, received_at);

