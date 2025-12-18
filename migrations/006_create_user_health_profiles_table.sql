-- Create USER_HEALTH_PROFILES table
CREATE TABLE IF NOT EXISTS user_health_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    age INTEGER CHECK (age > 0),
    gender VARCHAR(20),
    medical_conditions JSONB,
    allergies JSONB,
    medications JSONB,
    activity_preferences JSONB,
    health_goals JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id (already unique, but index helps with joins)
CREATE INDEX IF NOT EXISTS idx_user_health_profiles_user_id ON user_health_profiles(user_id);

-- Create trigger to auto-update updated_at
DROP TRIGGER IF EXISTS update_user_health_profiles_updated_at ON user_health_profiles;
CREATE TRIGGER update_user_health_profiles_updated_at BEFORE UPDATE ON user_health_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

