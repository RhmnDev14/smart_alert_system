-- Create ACTIVITY_COMPLETIONS table
CREATE TABLE IF NOT EXISTS activity_completions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on activity_id
CREATE INDEX IF NOT EXISTS idx_activity_completions_activity_id ON activity_completions(activity_id);
CREATE INDEX IF NOT EXISTS idx_activity_completions_is_completed ON activity_completions(is_completed);

