-- Seed data for development / demo
-- ============================================================

-- Hospitals
INSERT INTO hospitals (name, code) VALUES
    ('โรงพยาบาล A', 'HOSP_A'),
    ('โรงพยาบาล B', 'HOSP_B')
ON CONFLICT (code) DO NOTHING;

-- Patients for Hospital A (id = 1)
INSERT INTO patients (hospital_id, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en, last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender)
VALUES
    (1, 'สมชาย',  '',      'ใจดี',    'Somchai',  '',      'Jaidee',  '1990-01-15', 'HN-000001', '1234567890123', 'AA1234567', '0812345678', 'somchai.j@email.com',  'M'),
    (1, 'สมหญิง', '',      'รักสุข',  'Somying',  '',      'Raksuk',  '1985-06-20', 'HN-000002', '9876543210987', 'BB9876543', '0898765432', 'somying.r@email.com',  'F'),
    (1, 'จอห์น',  'เจมส์', 'สมิธ',    'John',     'James', 'Smith',   '1978-11-30', 'HN-000003', '',              'CC5551234', '0856789012', 'john.smith@email.com', 'M');

-- Patients for Hospital B (id = 2)
INSERT INTO patients (hospital_id, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en, last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender)
VALUES
    (2, 'มานะ',   '',      'ทำดี',    'Mana',     '',      'Thamdee', '1992-03-10', 'HN-B00001', '1111111111111', 'DD1111111', '0811111111', 'mana.t@email.com',     'M'),
    (2, 'วิไล',   '',      'สุขใจ',   'Wilai',    '',      'Sukjai',  '1988-12-05', 'HN-B00002', '2222222222222', 'EE2222222', '0822222222', 'wilai.s@email.com',    'F');
