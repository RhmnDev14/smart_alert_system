-- Create SCHEDULED_ALERTS table
CREATE TABLE IF NOT EXISTS scheduled_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alert_type VARCHAR(20) NOT NULL CHECK (alert_type IN ('morning', 'evening')),
    alert_time TIME NOT NULL,
    is_active BOOLEAN DEFAULT true,
    template TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_scheduled_alerts_user_id ON scheduled_alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_alerts_alert_type ON scheduled_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_scheduled_alerts_is_active ON scheduled_alerts(is_active);
CREATE INDEX IF NOT EXISTS idx_scheduled_alerts_user_alert_type ON scheduled_alerts(user_id, alert_type);

-- Create trigger to auto-update updated_at
DROP TRIGGER IF EXISTS update_scheduled_alerts_updated_at ON scheduled_alerts;
CREATE TRIGGER update_scheduled_alerts_updated_at BEFORE UPDATE ON scheduled_alerts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

