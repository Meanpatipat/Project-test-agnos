-- Hospital Middleware Database Schema
-- PostgreSQL 16

-- ============================================================
-- Table: hospitals
-- ============================================================
CREATE TABLE IF NOT EXISTS hospitals (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(50)  NOT NULL UNIQUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE  hospitals       IS 'Registered hospitals in the middleware system';
COMMENT ON COLUMN hospitals.code  IS 'Unique short code for the hospital (e.g. HOSP_A)';

-- ============================================================
-- Table: patients
-- Each patient belongs to exactly one hospital.
-- ============================================================
CREATE TABLE IF NOT EXISTS patients (
    id              SERIAL PRIMARY KEY,
    hospital_id     INTEGER       NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    first_name_th   VARCHAR(255)  DEFAULT '',
    middle_name_th  VARCHAR(255)  DEFAULT '',
    last_name_th    VARCHAR(255)  DEFAULT '',
    first_name_en   VARCHAR(255)  DEFAULT '',
    middle_name_en  VARCHAR(255)  DEFAULT '',
    last_name_en    VARCHAR(255)  DEFAULT '',
    date_of_birth   DATE,
    patient_hn      VARCHAR(50)   NOT NULL,
    national_id     VARCHAR(13)   DEFAULT '',
    passport_id     VARCHAR(20)   DEFAULT '',
    phone_number    VARCHAR(20)   DEFAULT '',
    email           VARCHAR(255)  DEFAULT '',
    gender          CHAR(1)       CHECK (gender IN ('M', 'F')),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE  patients             IS 'Patient records scoped to a hospital';
COMMENT ON COLUMN patients.hospital_id IS 'FK to the hospital this patient belongs to';
COMMENT ON COLUMN patients.patient_hn  IS 'Hospital Number – unique within the hospital';
COMMENT ON COLUMN patients.gender      IS 'M = Male, F = Female';

-- Unique HN per hospital
CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_hospital_hn
    ON patients (hospital_id, patient_hn);

-- Lookup indexes for search fields
CREATE INDEX IF NOT EXISTS idx_patients_national_id   ON patients (national_id)   WHERE national_id  <> '';
CREATE INDEX IF NOT EXISTS idx_patients_passport_id   ON patients (passport_id)   WHERE passport_id  <> '';
CREATE INDEX IF NOT EXISTS idx_patients_name_en       ON patients (first_name_en, last_name_en);
CREATE INDEX IF NOT EXISTS idx_patients_phone         ON patients (phone_number)  WHERE phone_number <> '';
CREATE INDEX IF NOT EXISTS idx_patients_email         ON patients (email)         WHERE email        <> '';
CREATE INDEX IF NOT EXISTS idx_patients_hospital      ON patients (hospital_id);

-- ============================================================
-- Table: staff
-- Each staff member belongs to exactly one hospital.
-- They can only query patients of that same hospital.
-- ============================================================
CREATE TABLE IF NOT EXISTS staff (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(100)  NOT NULL,
    password_hash VARCHAR(255)  NOT NULL,
    hospital_id   INTEGER       NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMENT ON TABLE  staff             IS 'Hospital staff with login credentials';
COMMENT ON COLUMN staff.hospital_id IS 'FK to the hospital – limits patient search scope';

-- Username is unique per hospital (same person can exist in different hospitals)
CREATE UNIQUE INDEX IF NOT EXISTS idx_staff_username_hospital
    ON staff (username, hospital_id);
