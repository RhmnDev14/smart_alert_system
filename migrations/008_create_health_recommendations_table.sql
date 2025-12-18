-- Create HEALTH_RECOMMENDATIONS table
CREATE TABLE IF NOT EXISTS health_recommendations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recommendation_type_id UUID REFERENCES recommendation_types(id) ON DELETE SET NULL,
    recommendation_text TEXT NOT NULL,
    activity_id UUID REFERENCES activities(id) ON DELETE SET NULL,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE,
    is_read BOOLEAN DEFAULT false,
    priority INTEGER DEFAULT 3 CHECK (priority >= 1 AND priority <= 5)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_health_recommendations_user_id ON health_recommendations(user_id);
CREATE INDEX IF NOT EXISTS idx_health_recommendations_recommendation_type_id ON health_recommendations(recommendation_type_id);
CREATE INDEX IF NOT EXISTS idx_health_recommendations_activity_id ON health_recommendations(activity_id);
CREATE INDEX IF NOT EXISTS idx_health_recommendations_is_read ON health_recommendations(is_read);
CREATE INDEX IF NOT EXISTS idx_health_recommendations_generated_at ON health_recommendations(generated_at);

