# Hospital Middleware API

ระบบ Middleware สำหรับโรงพยาบาล พัฒนาด้วย Go, Gin Framework, PostgreSQL, Docker และ Nginx

## Tech Stack

- **Go** + **Gin Framework** — API server
- **PostgreSQL 16** — Database
- **Docker & Docker Compose** — Containerization
- **Nginx** — Reverse proxy with rate limiting
- **JWT** — Authentication
- **bcrypt** — Password hashing
- **GORM** — ORM

## Project Structure

```
hospital-middleware/
├── main.go                          # Entry point
├── Dockerfile                       # Multi-stage Docker build
├── docker-compose.yml               # PostgreSQL + App + Nginx
├── .env                             # Local dev environment
├── config/
│   └── config.go                    # Configuration management
├── database/
│   └── database.go                  # PostgreSQL connection
├── migrations/
│   ├── init.sql                     # Database schema
│   └── seed.sql                     # Sample data
├── models/
│   ├── patient.go                   # Patient model + search request
│   ├── staff.go                     # Staff model + DTOs
│   ├── hospital.go                  # Hospital model
│   └── response.go                  # API response wrappers
├── repository/
│   ├── patient_repository.go        # Repository interfaces
│   ├── postgres_repository.go       # PostgreSQL implementations
│   └── mock_patient_repository.go   # Mock implementations (tests)
├── middleware/
│   └── auth.go                      # JWT authentication
├── handler/
│   ├── staff_handler.go             # Staff create + login
│   ├── patient_handler.go           # Patient search
│   ├── staff_handler_test.go        # Staff tests (11 cases)
│   └── patient_handler_test.go      # Patient tests (12 cases)
├── router/
│   └── router.go                    # Route configuration
└── nginx/
    └── nginx.conf                   # Reverse proxy config
```

## Quick Start

### Option 1: Docker Compose (recommended)
```bash
docker compose up --build
```
Access via: `http://localhost` (Nginx) or `http://localhost:8080` (direct)

### Option 2: Local Development
```bash
# Start PostgreSQL only
docker compose up postgres -d

# Run the app
go run main.go
```

### Run Tests
```bash
go test ./... -v
```

## API Endpoints

### 1. Health Check
```
GET /health
```

### 2. Create Staff (`/staff/create`)
```
POST /staff/create
Content-Type: application/json

{
  "username": "nurse_a",
  "password": "secret123",
  "hospital": "HOSP_A"
}
```

**Responses:**
- `201` — Staff created
- `400` — Invalid input / hospital not found
- `409` — Username already exists in this hospital

### 3. Staff Login (`/staff/login`)
```
POST /staff/login
Content-Type: application/json

{
  "username": "nurse_a",
  "password": "secret123",
  "hospital": "HOSP_A"
}
```

**Success (200):**
```json
{
  "status": 200,
  "message": "Login successful.",
  "data": {
    "token": "eyJhbGci...",
    "username": "nurse_a",
    "hospital": "HOSP_A"
  }
}
```

### 4. Search Patient (`/patient/search`) — Requires Login
```
GET /patient/search?national_id=1234567890123
Authorization: Bearer <token>
```

**Optional query parameters (all fields optional):**

| Parameter     | Description                |
|---------------|----------------------------|
| national_id   | National ID                |
| passport_id   | Passport ID                |
| first_name    | First name (TH or EN)      |
| middle_name   | Middle name (TH or EN)     |
| last_name     | Last name (TH or EN)       |
| date_of_birth | Date of birth (YYYY-MM-DD) |
| phone_number  | Phone number               |
| email         | Email address              |

> **Hospital isolation:** Staff can only see patients belonging to their own hospital.

## Database Schema

```
hospitals (id, name, code, created_at, updated_at)
patients  (id, hospital_id FK, names, dob, hn, ids, contact, gender, timestamps)
staff     (id, username, password_hash, hospital_id FK, timestamps)
```

- `patients.hospital_id` → scopes data per hospital
- `staff(username, hospital_id)` → unique constraint
- `patients(hospital_id, patient_hn)` → unique HN per hospital

## Test Coverage (23 tests)

| Test | Status |
|------|--------|
| **Staff Create** | |
| ✅ Success | PASS |
| ✅ Duplicate username | PASS |
| ✅ Invalid hospital | PASS |
| ✅ Missing fields | PASS |
| ✅ Short password | PASS |
| ✅ Short username | PASS |
| **Staff Login** | |
| ✅ Success + token | PASS |
| ✅ Wrong password | PASS |
| ✅ Non-existent user | PASS |
| ✅ Wrong hospital | PASS |
| ✅ Missing fields | PASS |
| **Patient Search** | |
| ✅ No filters (all) | PASS |
| ✅ By national_id | PASS |
| ✅ By passport_id | PASS |
| ✅ By first_name | PASS |
| ✅ By email | PASS |
| ✅ By date_of_birth | PASS |
| ✅ No auth token | PASS |
| ✅ Invalid token | PASS |
| ✅ Hospital isolation | PASS |
| ✅ Hospital B own data | PASS |
| ✅ No match | PASS |
| ✅ Health check | PASS |

## Configuration

| Variable            | Default                      | Description          |
|---------------------|------------------------------|----------------------|
| PORT                | 8080                         | Server port          |
| DB_HOST             | localhost                    | PostgreSQL host      |
| DB_PORT             | 5432                         | PostgreSQL port      |
| DB_USER             | hospital_admin               | DB username          |
| DB_PASSWORD         | hospital_secret_2026         | DB password          |
| DB_NAME             | hospital_middleware           | DB name              |
| DB_SSLMODE          | disable                      | SSL mode             |
| JWT_SECRET          | hospital-middleware-jwt-...   | JWT signing key      |
