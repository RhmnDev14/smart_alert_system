-- Create ACTIVITY_CATEGORIES table
CREATE TABLE IF NOT EXISTS activity_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on name
CREATE INDEX IF NOT EXISTS idx_activity_categories_name ON activity_categories(name);

