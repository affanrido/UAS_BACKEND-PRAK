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
- ✅ File validation & security

### FR-004: Submit untuk Verifikasi ✅
- ✅ Mahasiswa submit prestasi draft
- ✅ Update status menjadi 'submitted'
- ✅ Create notification untuk dosen wali
- ✅ Return updated status
- ✅ In-app notification system

### FR-005: Hapus Prestasi ✅
- ✅ Soft delete data di MongoDB
- ✅ Update reference di PostgreSQL
- ✅ Only draft status can be deleted
- ✅ Ownership validation
- ✅ Return success message

### FR-006: View Prestasi Mahasiswa Bimbingan ✅
- ✅ Get list student IDs dari tabel students where advisor_id
- ✅ Get achievements references dengan filter student_ids
- ✅ Fetch detail dari MongoDB
- ✅ Return list dengan pagination
- ✅ Combined data (reference + achievement + student info)

### FR-007: Verify Prestasi ✅
- ✅ Dosen review prestasi detail
- ✅ Dosen approve/reject prestasi
- ✅ Update status menjadi 'verified' atau 'rejected'
- ✅ Set verified_by dan verified_at
- ✅ Return updated status
- ✅ Notification untuk mahasiswa

### FR-008: Reject Prestasi ✅
- ✅ Dosen input rejection note
- ✅ Update status menjadi 'rejected'
- ✅ Save rejection_note
- ✅ Create notification untuk mahasiswa
- ✅ Return updated status
- ✅ **Implemented in FR-007** (approved: false)

### FR-009: Manage Users ✅
- ✅ Create/update/delete user
- ✅ Assign role
- ✅ Set student/lecturer profile
- ✅ Set advisor untuk mahasiswa
- ✅ Full CRUD operations
- ✅ Admin-only access

### FR-010: View All Achievements ✅
- ✅ Get all achievement references
- ✅ Fetch details dari MongoDB
- ✅ Apply filters dan sorting
- ✅ Return dengan pagination
- ✅ Summary statistics
- ✅ Combined data response

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

#### Mengapa Ada File .sql? 🗄️

File `.sql` dalam folder `database/` adalah **script database** yang diperlukan untuk setup sistem:

- **`schema.sql`** 🏗️ - Membuat struktur database (tabel, relasi, index)
- **`seed.sql`** 🌱 - Mengisi data awal untuk testing (roles, permissions, test users)

#### Setup PostgreSQL Database:

```bash
# Buat database PostgreSQL
createdb uas_backend

# Jalankan schema untuk membuat struktur tabel
psql -U postgres -d uas_backend -f database/schema.sql

# Jalankan seed untuk mengisi data testing
psql -U postgres -d uas_backend -f database/seed.sql

# Atau gunakan Makefile (lebih mudah)
make db-setup
```

#### Setup MongoDB:
MongoDB akan otomatis membuat database dan collection saat aplikasi pertama kali berjalan.

#### Reset Database (jika diperlukan):
```bash
# Reset complete database
make db-reset

# Atau manual
dropdb uas_backend
createdb uas_backend
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

### 🆕 API v1 Endpoints (New Implementation)

Implementasi API baru dengan versioning dan struktur yang lebih terorganisir:

#### **Base URL**: `/api/v1/`

#### **5.1 Authentication**
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/profile` - Get user profile

#### **5.2 Users (Admin)**
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `PUT /api/v1/users/:id/role` - Assign role

#### **5.4 Achievements**
- `GET /api/v1/achievements` - List achievements (filtered by role)
- `GET /api/v1/achievements/:id` - Get achievement detail
- `POST /api/v1/achievements` - Create achievement (Mahasiswa)
- `PUT /api/v1/achievements/:id` - Update achievement (Mahasiswa)
- `DELETE /api/v1/achievements/:id` - Delete achievement (Mahasiswa)
- `POST /api/v1/achievements/:id/submit` - Submit for verification
- `POST /api/v1/achievements/:id/verify` - Verify achievement (Dosen Wali)
- `POST /api/v1/achievements/:id/reject` - Reject achievement (Dosen Wali)
- `GET /api/v1/achievements/:id/history` - Status history
- `POST /api/v1/achievements/:id/attachments` - Upload files

#### **5.5 Students & Lecturers**
- `GET /api/v1/students` - Get all students
- `GET /api/v1/students/:id` - Get student by ID
- `GET /api/v1/students/:id/achievements` - Get student achievements
- `PUT /api/v1/students/:id/advisor` - Set student advisor
- `GET /api/v1/lecturers` - Get all lecturers
- `GET /api/v1/lecturers/:id/advisees` - Get lecturer advisees

#### **5.8 Reports & Analytics**
- `GET /api/v1/reports/statistics` - Get role-based statistics
- `GET /api/v1/reports/student/:id` - Get student report

**Features:**
- ✅ **Consistent Response Format** - Standardized JSON responses
- ✅ **Role-Based Filtering** - Automatic filtering based on user role
- ✅ **Comprehensive Pagination** - Page, limit, total, total_pages
- ✅ **Advanced Filtering** - Multiple query parameters for filtering
- ✅ **File Upload Support** - Multipart form data handling
- ✅ **Detailed Error Responses** - Structured error messages
- ✅ **Status History Tracking** - Achievement status changes
- ✅ **Analytics & Reports** - Role-based statistics and reports

See [docs/API_V1_ENDPOINTS.md](docs/API_V1_ENDPOINTS.md) for detailed documentation.

---

### Legacy API Endpoints (Existing)

#### Public Endpoints

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
UAS_BACKEND/
├── database/                       # 🗄️ Database setup files
│   ├── schema.sql                  # PostgreSQL database structure
│   ├── seed.sql                    # Test data untuk development
│   └── README.md                   # Database documentation
├── docs/                           # 📚 API documentation
│   ├── swagger.yaml                # OpenAPI 3.0 specification
│   ├── swagger.json                # JSON format for tools
│   ├── index.html                  # Custom Swagger UI
│   ├── README.md                   # API documentation guide
│   └── *.postman_*                 # Postman collection & environment
├── domain/                         # 🏗️ Business logic layer
│   ├── config/                     # Configuration management
│   │   ├── config.go               # Database configuration
│   │   ├── env.go                  # Environment variables
│   │   ├── logger.go               # Logging configuration
│   │   ├── mongodb.go              # MongoDB configuration
│   │   └── token.go                # JWT token configuration
│   ├── middleware/                 # 🛡️ Authentication & RBAC
│   │   ├── Auth.go                 # JWT authentication middleware
│   │   ├── rbac.go                 # Role-based access control
│   │   └── TokenMiddleware.go      # Token validation middleware
│   ├── model/                      # 📊 Data models (PostgreSQL + MongoDB)
│   │   ├── Users.go                # User model
│   │   ├── Roles.go                # Role model
│   │   ├── Permission.go           # Permission model
│   │   ├── Role_Permission.go      # Role-Permission mapping
│   │   ├── Student.go              # Student profile model
│   │   ├── Lecturers.go            # Lecturer profile model
│   │   ├── achievement_references.go # Achievement reference (PostgreSQL)
│   │   ├── achievements.go         # Achievement detail (MongoDB)
│   │   ├── Notification.go         # Notification model
│   │   ├── Pagination.go           # Pagination helper
│   │   └── README.md               # Model documentation
│   ├── repository/                 # 🗃️ Database access layer
│   │   ├── authRepo.go             # Authentication repository
│   │   ├── achievementRepo.go      # Achievement repository
│   │   ├── userRepo.go             # User management repository
│   │   ├── statisticsRepo.go       # Statistics repository
│   │   ├── notificationRepo.go     # Notification repository
│   │   └── rbacRepo.go             # RBAC repository
│   ├── route/                      # 🌐 HTTP handlers & routing
│   │   ├── authRoute.go            # Authentication routes
│   │   ├── achievementRoute.go     # Achievement routes
│   │   ├── lecturerRoute.go        # Lecturer routes
│   │   ├── adminRoute.go           # Admin routes
│   │   ├── statisticsRoute.go      # Statistics routes
│   │   └── notificationRoute.go    # Notification routes
│   └── service/                    # 🔧 Business logic services
│       ├── Authservice.go          # Authentication service
│       ├── achievementService.go   # Achievement service
│       ├── adminAchievementService.go # Admin achievement service
│       ├── userService.go          # User management service
│       ├── statisticsService.go    # Statistics service
│       ├── notificationService.go  # Notification service
│       ├── rbacservice.go          # RBAC service
│       └── fileService.go          # File upload service
├── tests/                          # 🧪 Comprehensive testing suite
│   ├── run_tests.go                # Test runner script
│   ├── test_swagger.go             # Swagger documentation tester
│   ├── test_config.go              # Common test configuration
│   ├── README.md                   # Testing documentation
│   ├── mocks/                      # 🎭 Mock implementations
│   │   ├── auth_repository_mock.go # Auth repository mock
│   │   ├── achievement_repository_mock.go # Achievement repository mock
│   │   ├── user_repository_mock.go # User repository mock
│   │   ├── statistics_repository_mock.go # Statistics repository mock
│   │   ├── notification_repository_mock.go # Notification repository mock
│   │   ├── auth_service_mock.go    # Auth service mock
│   │   └── notification_service_mock.go # Notification service mock
│   ├── service/                    # 🔧 Service layer tests
│   │   ├── auth_service_test.go    # Authentication service tests
│   │   ├── achievement_service_test.go # Achievement service tests
│   │   ├── user_service_test.go    # User management service tests
│   │   ├── statistics_service_test.go # Statistics service tests
│   │   └── notification_service_test.go # Notification service tests
│   ├── middleware/                 # 🛡️ Middleware tests
│   │   ├── auth_middleware_test.go # Authentication middleware tests
│   │   └── rbac_middleware_test.go # RBAC middleware tests
│   └── integration/                # 🔗 Integration tests
│       └── auth_handler_test.go    # HTTP handler integration tests
├── tools/                          # 🔧 Utility tools
│   └── generate_hash.go            # Password hash generator
├── .env.example                    # Environment variables template
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
├── main.go                         # Application entry point
├── Makefile                        # Build automation commands
├── DEPLOYMENT.md                   # Production deployment guide
└── README.md                       # Project documentation (this file)
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

## Testing

### Unit Testing ✅

Comprehensive unit testing implementation following requirements:

- ✅ **Test individual functions and methods**
- ✅ **Mock external dependencies** 
- ✅ **Cover success cases, error cases, and edge cases**
- ✅ **Use testify framework for assertions and mocks**

#### Test Structure

```
tests/
├── mocks/              # Mock implementations
├── service/            # Service layer tests  
├── middleware/         # Middleware tests
├── integration/        # Integration tests
└── README.md          # Testing documentation
```

#### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage report
make test-coverage

# Run specific test suites
make test-service      # Service layer only
make test-middleware   # Middleware only
make test-integration  # Integration tests only

# Or using go test directly
go test ./tests/... -v
go test ./tests/service -v
go test ./tests/... -v -cover
```

#### Test Coverage

- **Service Layer**: AuthService, AchievementService, UserService, StatisticsService, NotificationService
- **Middleware**: Authentication, RBAC (Role-Based Access Control)
- **Integration**: HTTP handlers, request/response validation
- **Mocks**: All external dependencies (repositories, services)

#### Test Features

- **Comprehensive Coverage**: Success cases, error cases, edge cases, validation errors
- **Mock Dependencies**: Database operations, external services, cross-service calls
- **Integration Testing**: HTTP endpoints, middleware chains, request validation
- **Performance Testing**: No external dependencies, fast execution
- **CI/CD Ready**: Deterministic results, clear pass/fail indicators

See [tests/README.md](tests/README.md) for detailed testing documentation.

## API Documentation

### 7.1 Swagger untuk API Documentation ✅

Comprehensive API documentation dengan Swagger UI yang interaktif:

#### **Access Documentation:**
- **Swagger UI**: `http://localhost:8080/swagger/` - Interactive API testing
- **Static Docs**: [docs/swagger.yaml](docs/swagger.yaml) - OpenAPI 3.0 specification
- **Postman Collection**: [docs/UAS_Backend_API.postman_collection.json](docs/UAS_Backend_API.postman_collection.json)

#### **Features:**
- ✅ **Interactive Testing** - Test endpoints directly from browser
- ✅ **Authentication Support** - JWT Bearer token integration
- ✅ **Request/Response Examples** - Complete examples for all endpoints
- ✅ **Schema Validation** - Request/response schema documentation
- ✅ **Error Documentation** - Comprehensive error codes and messages
- ✅ **Role-Based Examples** - Different examples for Admin/Lecturer/Student

#### **Quick Start:**
```bash
# Start server
go run main.go

# Open Swagger UI
open http://localhost:8080/swagger/

# Test login endpoint
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"admin@example.com","password":"password123"}'
```

#### **Documentation Structure:**
```
docs/
├── swagger.yaml              # OpenAPI 3.0 specification
├── swagger.json             # JSON format for tools
├── index.html               # Custom Swagger UI
├── README.md                # API documentation guide
├── UAS_Backend_API.postman_collection.json    # Postman collection
└── UAS_Backend_Environment.postman_environment.json # Postman environment
```

#### **API Coverage:**
- **Authentication**: Login, JWT validation
- **Achievements**: Submit, verify, delete, list (FR-003 to FR-008, FR-010)
- **User Management**: CRUD operations (FR-009)
- **Statistics**: Role-based statistics (FR-011)
- **Notifications**: Real-time notification system
- **RBAC**: Role-based access control (FR-002)

#### **Testing Tools:**
- **Swagger UI**: Interactive browser-based testing
- **Postman**: Complete collection with environment variables
- **cURL**: Command-line examples for all endpoints
- **Automated Tests**: Unit tests with comprehensive coverage

See [docs/README.md](docs/README.md) for detailed API documentation.

### 7.2 Github Repository ✅

Repository structure dan documentation:

#### **Repository Features:**
- ✅ **Complete Source Code** - All functional requirements implemented
- ✅ **Comprehensive Documentation** - README, API docs, testing guides
- ✅ **Database Schema** - PostgreSQL + MongoDB setup scripts
- ✅ **Environment Configuration** - .env.example with all required variables
- ✅ **Testing Suite** - Unit tests, integration tests, mocks
- ✅ **API Documentation** - Swagger UI, Postman collection
- ✅ **Build Instructions** - Step-by-step setup guide

#### **Repository Structure:**
```
UAS_BACKEND/
├── domain/                  # Business logic layer
│   ├── config/             # Configuration management
│   ├── middleware/         # Authentication & RBAC
│   ├── model/              # Data models (PostgreSQL + MongoDB)
│   ├── repository/         # Database layer
│   ├── route/              # HTTP handlers & routing
│   └── service/            # Business logic services
├── database/               # Database setup
│   ├── schema.sql          # PostgreSQL schema
│   └── seed.sql            # Test data
├── docs/                   # API documentation
│   ├── swagger.yaml        # OpenAPI specification
│   ├── README.md           # API guide
│   └── *.postman_*         # Postman collection
├── tests/                  # Testing suite
│   ├── mocks/              # Mock implementations
│   ├── service/            # Service layer tests
│   ├── middleware/         # Middleware tests
│   └── integration/        # Integration tests
├── tools/                  # Utility tools
├── .env.example            # Environment template
├── go.mod                  # Go dependencies
├── main.go                 # Application entry point
├── Makefile               # Build automation
└── README.md              # Project documentation
```

#### **Documentation Quality:**
- **Setup Instructions** - Complete environment setup
- **API Reference** - All endpoints documented with examples
- **Testing Guide** - How to run and write tests
- **Architecture Overview** - System design and patterns
- **Deployment Guide** - Production deployment instructions

## License

MIT
