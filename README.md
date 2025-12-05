# UAS Backend - Authentication System

Backend API untuk sistem autentikasi dengan RBAC (Role-Based Access Control).

## Features

- ✅ Login dengan username/email dan password
- ✅ JWT token authentication
- ✅ Role-Based Access Control (RBAC)
- ✅ Password hashing dengan bcrypt
- ✅ Status user aktif/non-aktif
- ✅ Permissions management

## Tech Stack

- Go 1.24
- Fiber v2 (Web Framework)
- PostgreSQL (Database)
- JWT (Authentication)
- Bcrypt (Password Hashing)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE uas_backend;
```

Jalankan schema dan seed:

```bash
psql -U postgres -d uas_backend -f database/schema.sql
psql -U postgres -d uas_backend -f database/seed.sql
```

### 3. Environment Variables

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
copy .env.example .env
```

Edit `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=uas_backend
JWT_SECRET=your-secret-key
```

### 4. Run Server

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Authentication

#### POST /api/auth/login

Login dengan username/email dan password.

**Request Body:**

```json
{
  "identifier": "admin@example.com",
  "password": "password123"
}
```

**Success Response (200):**

```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "Administrator",
      "role_id": "550e8400-e29b-41d4-a716-446655440001",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

**Error Response (401):**

```json
{
  "error": "invalid credentials"
}
```

**Error Response (400):**

```json
{
  "error": "Username/email and password are required"
}
```

### Health Check

#### GET /health

Check server status.

**Response:**

```json
{
  "status": "ok"
}
```

## Test Users

Setelah menjalankan seed.sql, tersedia 3 test users:

| Username  | Email                 | Password    | Role     |
|-----------|-----------------------|-------------|----------|
| admin     | admin@example.com     | password123 | admin    |
| lecturer1 | lecturer@example.com  | password123 | lecturer |
| student1  | student@example.com   | password123 | student  |

## Testing dengan cURL

### Login sebagai Admin

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"admin@example.com\",\"password\":\"password123\"}"
```

### Login dengan Username

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"admin\",\"password\":\"password123\"}"
```

## Project Structure

```
.
├── domain/
│   ├── config/          # Configuration (DB, JWT)
│   ├── middleware/      # Auth & RBAC middleware
│   ├── model/           # Data models & DTOs
│   ├── repository/      # Database layer
│   ├── route/           # HTTP handlers & routes
│   └── service/         # Business logic
├── database/
│   ├── schema.sql       # Database schema
│   └── seed.sql         # Seed data
├── .env.example         # Environment variables template
├── go.mod
├── go.sum
└── main.go              # Application entry point
```

## FR-001: Login Flow

1. ✅ User mengirim kredensial (username/email + password)
2. ✅ Sistem memvalidasi kredensial dari database
3. ✅ Sistem mengecek status aktif user
4. ✅ Sistem generate JWT token dengan role dan permissions
5. ✅ Return token dan user profile

## License

MIT
