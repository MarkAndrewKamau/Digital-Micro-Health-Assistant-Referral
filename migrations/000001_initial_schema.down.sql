-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_appointments_updated_at ON appointments;
DROP TRIGGER IF EXISTS update_referrals_updated_at ON referrals;
DROP TRIGGER IF EXISTS update_triage_sessions_updated_at ON triage_sessions;
DROP TRIGGER IF EXISTS update_clinicians_updated_at ON clinicians;
DROP TRIGGER IF EXISTS update_facilities_updated_at ON facilities;
DROP TRIGGER IF EXISTS update_patients_updated_at ON patients;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS otp_codes;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS embeddings;
DROP TABLE IF EXISTS consent_logs;
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS referrals;
DROP TABLE IF EXISTS triage_sessions;
DROP TABLE IF EXISTS clinicians;
DROP TABLE IF EXISTS facilities;
DROP TABLE IF EXISTS patients;

-- Drop types
DROP TYPE IF EXISTS consent_type;
DROP TYPE IF EXISTS facility_type;
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS referral_status;
DROP TYPE IF EXISTS triage_level;
DROP TYPE IF EXISTS user_role;

-- Drop extensions
DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS postgis;
DROP EXTENSION IF EXISTS "uuid-ossp";