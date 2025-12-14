# API v1 Endpoints Documentation

## Base URL
- **Development**: `http://localhost:8080/api/v1`
- **Production**: `https://api.yourdomain.com/api/v1`

## Authentication
All protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "pagination": { ... },  // For list endpoints
  "metadata": { ... }     // Additional info
}
```

## 5.1 Authentication

### POST /api/v1/auth/login
Login with username/email and password.

**Request:**
```json
{
  "identifier": "admin@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "770e8400-e29b-41d4-a716-446655440001",
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "Administrator",
      "role_id": "550e8400-e29b-41d4-a716-446655440001",
      "is_active": true
    },
    "expires_at": "2024-01-02T00:00:00Z"
  }
}
```

### POST /api/v1/auth/refresh
Refresh JWT token using refresh token.

### POST /api/v1/auth/logout
Logout and invalidate token.

### GET /api/v1/auth/profile
Get current user profile information.

## 5.2 Users (Admin Only)

### GET /api/v1/users
Get all users with pagination and filters.

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10)
- `role` (string): Filter by role (admin, lecturer, student)
- `is_active` (boolean): Filter by active status

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "user": { ... },
      "student": { ... },
      "lecturer": { ... },
      "role": { ... }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### GET /api/v1/users/:id
Get specific user by ID.

### POST /api/v1/users
Create new user with role and profile.

### PUT /api/v1/users/:id
Update user information.

### DELETE /api/v1/users/:id
Delete user.

### PUT /api/v1/users/:id/role
Assign role to user.

## 5.4 Achievements

### GET /api/v1/achievements
List achievements (filtered by role).

**Query Parameters:**
- `page`, `limit`: Pagination
- `status`: Filter by status (draft, submitted, verified, rejected)
- `type`: Filter by achievement type
- `student_id`: Filter by specific student

### GET /api/v1/achievements/:id
Get achievement detail.

### POST /api/v1/achievements
Create new achievement (Mahasiswa only).

### PUT /api/v1/achievements/:id
Update achievement (Mahasiswa only).

### DELETE /api/v1/achievements/:id
Delete achievement (Mahasiswa only).

### POST /api/v1/achievements/:id/submit
Submit achievement for verification.

### POST /api/v1/achievements/:id/verify
Verify achievement (Dosen Wali only).

### POST /api/v1/achievements/:id/reject
Reject achievement (Dosen Wali only).

### GET /api/v1/achievements/:id/history
Get achievement status history.

### POST /api/v1/achievements/:id/attachments
Upload files to achievement.
## 5.5 Students & Lecturers

### GET /api/v1/students
Get all students with filters.

**Query Parameters:**
- `page`, `limit`: Pagination
- `program_study`: Filter by program study
- `academic_year`: Filter by academic year
- `advisor_id`: Filter by advisor

### GET /api/v1/students/:id
Get specific student by ID.

### GET /api/v1/students/:id/achievements
Get achievements for specific student.

### PUT /api/v1/students/:id/advisor
Set advisor for student (Admin only).

### GET /api/v1/lecturers
Get all lecturers.

### GET /api/v1/lecturers/:id/advisees
Get students under lecturer supervision.

## 5.8 Reports & Analytics

### GET /api/v1/reports/statistics
Get role-based statistics.

**Query Parameters:**
- `start_date`: Start date (YYYY-MM-DD)
- `end_date`: End date (YYYY-MM-DD)
- `type`: Report type (overview, detailed, trends)

**Response varies by role:**
- **Admin**: System-wide statistics
- **Lecturer**: Advisee statistics
- **Student**: Own statistics

### GET /api/v1/reports/student/:id
Get detailed report for specific student.

**Query Parameters:**
- `start_date`, `end_date`: Date range
- `include_details`: Include detailed breakdown

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "error": "Invalid JSON format"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "error": "Missing or invalid token"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "Internal server error"
}
```