-- Enable required extensions (skip vector for now)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE user_role AS ENUM ('patient', 'chv', 'clinician', 'admin');
CREATE TYPE triage_level AS ENUM ('red', 'yellow', 'green');
CREATE TYPE referral_status AS ENUM ('pending', 'accepted', 'completed', 'cancelled');
CREATE TYPE appointment_status AS ENUM ('scheduled', 'confirmed', 'completed', 'cancelled', 'no_show');
CREATE TYPE facility_type AS ENUM ('dispensary', 'health_center', 'sub_county_hospital', 'county_hospital', 'private_clinic');
CREATE TYPE consent_type AS ENUM ('data_collection', 'data_sharing', 'sms_notifications', 'research');

-- Patients table
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255),
    date_of_birth DATE,
    gender VARCHAR(10),
    preferred_language VARCHAR(10) DEFAULT 'en',
    consent_flags JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_patients_phone ON patients(phone);
CREATE INDEX idx_patients_consent_flags ON patients USING GIN(consent_flags);

-- Facilities table (without PostGIS - using simple lat/lng columns)
CREATE TABLE facilities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type facility_type NOT NULL,
    level INTEGER,
    county VARCHAR(100),
    sub_county VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    services JSONB DEFAULT '[]',
    operating_hours JSONB DEFAULT '{}',
    accepts_referrals BOOLEAN DEFAULT true,
    accepts_mpesa BOOLEAN DEFAULT false,
    bed_capacity INTEGER,
    staff_count INTEGER,
    available_slots JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_facilities_type ON facilities(type);
CREATE INDEX idx_facilities_county ON facilities(county);
CREATE INDEX idx_facilities_services ON facilities USING GIN(services);
CREATE INDEX idx_facilities_lat_lng ON facilities(latitude, longitude);

-- Clinicians table
CREATE TABLE clinicians (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255),
    facility_id UUID REFERENCES facilities(id) ON DELETE SET NULL,
    specialization VARCHAR(255),
    license_number VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clinicians_facility ON clinicians(facility_id);
CREATE INDEX idx_clinicians_phone ON clinicians(phone);

-- Triage sessions table
CREATE TABLE triage_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    symptoms JSONB NOT NULL,
    summary_text TEXT,
    triage_level triage_level,
    triage_code VARCHAR(10),
    confidence DECIMAL(3, 2),
    recommended_action TEXT,
    llm_response JSONB,
    channel VARCHAR(20) DEFAULT 'sms',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_triage_patient ON triage_sessions(patient_id);
CREATE INDEX idx_triage_level ON triage_sessions(triage_level);
CREATE INDEX idx_triage_created ON triage_sessions(created_at DESC);
CREATE INDEX idx_triage_symptoms ON triage_sessions USING GIN(symptoms);

-- Referrals table
CREATE TABLE referrals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    triage_session_id UUID REFERENCES triage_sessions(id) ON DELETE SET NULL,
    facility_id UUID REFERENCES facilities(id) ON DELETE CASCADE,
    referral_token VARCHAR(20) UNIQUE NOT NULL,
    status referral_status DEFAULT 'pending',
    priority triage_level,
    notes TEXT,
    created_by_chv UUID,
    accepted_by_clinician UUID REFERENCES clinicians(id) ON DELETE SET NULL,
    accepted_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_referrals_patient ON referrals(patient_id);
CREATE INDEX idx_referrals_facility ON referrals(facility_id);
CREATE INDEX idx_referrals_token ON referrals(referral_token);
CREATE INDEX idx_referrals_status ON referrals(status);
CREATE INDEX idx_referrals_created ON referrals(created_at DESC);

-- Appointments table
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    referral_id UUID REFERENCES referrals(id) ON DELETE CASCADE,
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    facility_id UUID REFERENCES facilities(id) ON DELETE CASCADE,
    clinician_id UUID REFERENCES clinicians(id) ON DELETE SET NULL,
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status appointment_status DEFAULT 'scheduled',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_appointments_referral ON appointments(referral_id);
CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_facility ON appointments(facility_id);
CREATE INDEX idx_appointments_scheduled ON appointments(scheduled_time);
CREATE INDEX idx_appointments_status ON appointments(status);

-- Consent logs table
CREATE TABLE consent_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    consent_type consent_type NOT NULL,
    granted BOOLEAN NOT NULL,
    details JSONB DEFAULT '{}',
    granted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_consent_patient ON consent_logs(patient_id);
CREATE INDEX idx_consent_type ON consent_logs(consent_type);
CREATE INDEX idx_consent_granted_at ON consent_logs(granted_at DESC);

-- Users table (for CHVs, clinicians, admins authentication)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    role user_role NOT NULL DEFAULT 'patient',
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    clinician_id UUID REFERENCES clinicians(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_patient ON users(patient_id);
CREATE INDEX idx_users_clinician ON users(clinician_id);

-- Sessions table (for session-based auth)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    user_agent TEXT,
    ip_address INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_token ON sessions(session_token);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers
CREATE TRIGGER update_patients_updated_at BEFORE UPDATE ON patients FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_facilities_updated_at BEFORE UPDATE ON facilities FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_clinicians_updated_at BEFORE UPDATE ON clinicians FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_triage_sessions_updated_at BEFORE UPDATE ON triage_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_referrals_updated_at BEFORE UPDATE ON referrals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_appointments_updated_at BEFORE UPDATE ON appointments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();