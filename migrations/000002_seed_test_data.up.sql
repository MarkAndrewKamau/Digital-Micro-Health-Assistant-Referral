-- Seed test facilities
INSERT INTO facilities (id, name, type, level, county, sub_county, location, phone, services, operating_hours, accepts_referrals) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Kamulu Health Center', 'health_center', 3, 'Machakos', 'Kangundo', ST_GeogFromText('POINT(37.1234 -1.2345)'), '+254712345678', '["outpatient", "maternity", "lab", "pharmacy"]', '{"monday": "08:00-17:00", "tuesday": "08:00-17:00", "wednesday": "08:00-17:00", "thursday": "08:00-17:00", "friday": "08:00-17:00"}', true),
('550e8400-e29b-41d4-a716-446655440002', 'Makueni County Hospital', 'county_hospital', 4, 'Makueni', 'Makueni', ST_GeogFromText('POINT(37.6234 -1.8345)'), '+254722345678', '["outpatient", "inpatient", "maternity", "lab", "pharmacy", "surgery", "icu"]', '{"monday": "24/7", "tuesday": "24/7", "wednesday": "24/7", "thursday": "24/7", "friday": "24/7", "saturday": "24/7", "sunday": "24/7"}', true);

-- Seed test clinicians
INSERT INTO clinicians (id, name, phone, email, facility_id, specialization) VALUES
('550e8400-e29b-41d4-a716-446655440010', 'Dr. Jane Wambui', '+254733345678', 'jane.wambui@health.go.ke', '550e8400-e29b-41d4-a716-446655440001', 'General Practice'),
('550e8400-e29b-41d4-a716-446655440011', 'Dr. John Omondi', '+254744345678', 'john.omondi@health.go.ke', '550e8400-e29b-41d4-a716-446655440002', 'Internal Medicine');

-- Seed test patient
INSERT INTO patients (id, phone, name, date_of_birth, gender, preferred_language, consent_flags) VALUES
('550e8400-e29b-41d4-a716-446655440020', '+254755345678', 'Test Patient', '1990-01-01', 'male', 'en', '{"data_collection": true, "sms_notifications": true}');

-- Seed test user (CHV)
INSERT INTO users (id, phone, name, role) VALUES
('550e8400-e29b-41d4-a716-446655440030', '+254766345678', 'Mary Muthoni (CHV)', 'chv');