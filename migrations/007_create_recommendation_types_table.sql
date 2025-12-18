-- Create RECOMMENDATION_TYPES table
CREATE TABLE IF NOT EXISTS recommendation_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    trigger_condition TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on name
CREATE INDEX IF NOT EXISTS idx_recommendation_types_name ON recommendation_types(name);

