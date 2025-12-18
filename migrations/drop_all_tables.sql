-- Script untuk drop semua tabel (HATI-HATI: Hanya untuk development!)
-- Jangan jalankan di production!

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS scheduled_alerts CASCADE;
DROP TABLE IF EXISTS alert_logs CASCADE;
DROP TABLE IF EXISTS message_history CASCADE;
DROP TABLE IF EXISTS health_recommendations CASCADE;
DROP TABLE IF EXISTS recommendation_types CASCADE;
DROP TABLE IF EXISTS user_health_profiles CASCADE;
DROP TABLE IF EXISTS activity_completions CASCADE;
DROP TABLE IF EXISTS activities CASCADE;
DROP TABLE IF EXISTS activity_categories CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Note: UUID extension tidak di-drop karena mungkin digunakan oleh database lain

