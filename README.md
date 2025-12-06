# UAS Backend - Authentication System

Backend API untuk sistem autentikasi dengan RBAC (Role-Based Access Control).

## Features

### FR-001: Login ✅
- ✅ Login dengan username/email dan password
- ✅ JWT token authentication
- ✅ Password hashing dengan bcrypt
- ✅ Status user aktif/non-aktif
- ✅ Generate JWT dengan role dan permissions

### FR-002: RBAC Middleware ✅
- ✅ JWT extraction dari header
- ✅ Token validation
- ✅ Permission check per endpoint
- ✅ In-memory cache dengan TTL
- ✅ Multiple permission strategies (any, all, single)
- ✅ Role-based access control

### FR-003: Submit Prestasi ✅
- ✅ Mahasiswa mengisi data prestasi
- ✅ Upload dokumen pendukung (max 10MB)
- ✅ Hybrid database (MongoDB + PostgreSQL)
- ✅ Status awal: 'draft'
- ✅ Multiple achievement types (competition, publication, organization, certification, academic, other)
- ✅ Dynamic details per type
- ✅ File validation & securit

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

### Public Endpoints

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

#### GET /health

Check server status.

**Response:**

```json
{
  "status": "ok"
}
```

---

### Protected Endpoints (Require Authentication)

All protected endpoints require `Authorization: Bearer <token>` header.

#### GET /api/profile

Get user profile (authentication only, no permission check).

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Success Response (200):**
```json
{
  "message": "Profile accessed",
  "user_id": "770e8400-e29b-41d4-a716-446655440001"
}
```

---

#### GET /api/students

Get students list. Requires `student.read` permission.

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Success Response (200):**
```json
{
  "message": "Students list",
  "permissions": ["student.read", "student.write"]
}
```

**Error Response (403):**
```json
{
  "error": "Forbidden: Insufficient permissions",
  "required": "student.read"
}
```

---

#### POST /api/students

Create new student. Requires `student.write` permission.

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

---

#### POST /api/achievements/:id/verify

Verify achievement. Requires `achievement.verify` permission.

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

---

#### GET /api/admin/dashboard

Admin dashboard. Requires `admin` role.

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

---

#### GET /api/reports

Reports endpoint. Requires any of: `student.read`, `lecturer.read`, or `achievement.read`.

**Headers:**
```
Authorization: Bearer <your_jwt_token>
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

## Feature Implementation

### FR-001: Login Flow ✅

1. ✅ User mengirim kredensial (username/email + password)
2. ✅ Sistem memvalidasi kredensial dari database
3. ✅ Sistem mengecek status aktif user
4. ✅ Sistem generate JWT token dengan role dan permissions
5. ✅ Return token dan user profile

### FR-002: RBAC Middleware Flow ✅

1. ✅ Ekstrak JWT dari header Authorization
2. ✅ Validasi token (signature, expiry)
3. ✅ Load user permissions dari cache/database
4. ✅ Check apakah user memiliki permission yang diperlukan
5. ✅ Allow/deny request

**Middleware Types:**
- `Authenticate()` - JWT validation only
- `RequirePermission(perm)` - Single permission check
- `RequireAnyPermission(perms...)` - Any of permissions (OR)
- `RequireAllPermissions(perms...)` - All permissions (AND)
- `RequireRole(role)` - Role-based check

**Caching:**
- In-memory cache with 5-minute TTL
- Auto cleanup expired entries
- Per-user permission cache

## License

MIT
