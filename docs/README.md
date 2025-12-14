# UAS Backend API Documentation

Comprehensive API documentation untuk UAS Backend system dengan authentication, RBAC, dan manajemen prestasi mahasiswa.

## 📚 Documentation Access

### Swagger UI (Interactive)
- **URL**: `http://localhost:8080/swagger/`
- **Features**: Interactive API testing, request/response examples, authentication support
- **Best for**: Development, testing, dan exploration

### Static Documentation
- **Swagger YAML**: [swagger.yaml](./swagger.yaml)
- **Swagger JSON**: [swagger.json](./swagger.json)
- **Best for**: Integration dengan tools lain, CI/CD

### Postman Collection
- **Collection**: [UAS_Backend_API.postman_collection.json](./UAS_Backend_API.postman_collection.json)
- **Environment**: [UAS_Backend_Environment.postman_environment.json](./UAS_Backend_Environment.postman_environment.json)
- **Best for**: Manual testing, team collaboration

## 🚀 Quick Start

### 1. Start the Server
```bash
go run main.go
```

### 2. Access Documentation
- Open browser: `http://localhost:8080/swagger/`
- Or use Postman with provided collection

### 3. Authentication Flow
1. **Login** dengan test credentials:
   - Admin: `admin@example.com` / `password123`
   - Lecturer: `lecturer@example.com` / `password123`
   - Student: `student@example.com` / `password123`

2. **Copy JWT token** dari response
3. **Add to Authorization header**: `Bearer <your_token>`

## 📋 API Overview

### Base URL
- **Development**: `http://localhost:8080`
- **Production**: `https://api.example.com`

### Authentication
```http
Authorization: Bearer <jwt_token>
```

### Response Format
```json
{
  "message": "Success message",
  "data": { ... },
  "pagination": { ... }  // For list endpoints
}
```

### Error Format
```json
{
  "error": "Error message",
  "details": "Additional details"  // Optional
}
```

## 🔐 Authentication & Authorization

### Login Endpoint
```http
POST /api/auth/login
Content-Type: application/json

{
  "identifier": "admin@example.com",
  "password": "password123"
}
```

### Response
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
      "is_active": true
    },
    "expires_at": "2024-01-02T00:00:00Z"
  }
}
```

### Role-Based Access Control (RBAC)

| Role | Permissions | Access |
|------|-------------|---------|
| **Admin** | All permissions | Full system access |
| **Lecturer** | `achievement.verify`, `student.read` | Verify student achievements |
| **Student** | `achievement.submit`, `achievement.read` | Submit own achievements |

## 📊 API Endpoints Summary

### Authentication
- `POST /api/auth/login` - User login

### Profile
- `GET /api/profile` - Get current user profile

### Achievements (Student)
- `GET /api/student/achievements` - Get own achievements
- `POST /api/student/achievements` - Submit new achievement
- `POST /api/student/achievements/{id}/submit` - Submit for verification
- `DELETE /api/student/achievements/{id}` - Delete draft achievement

### Achievements (Lecturer)
- `GET /api/lecturer/achievements` - Get advisee achievements
- `POST /api/lecturer/achievements/{id}/verify` - Verify/reject achievement

### Achievements (Admin)
- `GET /api/admin/achievements` - Get all achievements with filters

### User Management (Admin)
- `GET /api/admin/users` - Get all users
- `POST /api/admin/users` - Create new user
- `GET /api/admin/users/{id}` - Get user by ID
- `PUT /api/admin/users/{id}` - Update user
- `DELETE /api/admin/users/{id}` - Delete user

### Statistics
- `GET /api/student/statistics` - Student achievement statistics
- `GET /api/lecturer/statistics` - Lecturer advisee statistics
- `GET /api/admin/statistics` - System-wide statistics

### Notifications
- `GET /api/notifications` - Get user notifications
- `GET /api/notifications/unread-count` - Get unread count
- `POST /api/notifications/{id}/read` - Mark as read
- `POST /api/notifications/read-all` - Mark all as read

### Health Check
- `GET /health` - Server health status

## 🎯 Functional Requirements Mapping

| FR | Feature | Endpoints | Status |
|----|---------|-----------|---------|
| FR-001 | Login | `POST /api/auth/login` | ✅ |
| FR-002 | RBAC Middleware | All protected endpoints | ✅ |
| FR-003 | Submit Prestasi | `POST /api/student/achievements` | ✅ |
| FR-004 | Submit untuk Verifikasi | `POST /api/student/achievements/{id}/submit` | ✅ |
| FR-005 | Hapus Prestasi | `DELETE /api/student/achievements/{id}` | ✅ |
| FR-006 | View Prestasi Mahasiswa Bimbingan | `GET /api/lecturer/achievements` | ✅ |
| FR-007 | Verify Prestasi | `POST /api/lecturer/achievements/{id}/verify` | ✅ |
| FR-008 | Reject Prestasi | `POST /api/lecturer/achievements/{id}/verify` (approved: false) | ✅ |
| FR-009 | Manage Users | `GET/POST/PUT/DELETE /api/admin/users` | ✅ |
| FR-010 | View All Achievements | `GET /api/admin/achievements` | ✅ |
| FR-011 | Achievement Statistics | `GET /api/{role}/statistics` | ✅ |

## 🧪 Testing

### Using Swagger UI
1. Open `http://localhost:8080/swagger/`
2. Click "Authorize" button
3. Enter: `Bearer <your_token>`
4. Test endpoints interactively

### Using Postman
1. Import collection: `UAS_Backend_API.postman_collection.json`
2. Import environment: `UAS_Backend_Environment.postman_environment.json`
3. Run "Login - Admin" to get token
4. Test other endpoints

### Using cURL
```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"admin@example.com","password":"password123"}'

# Get profile (replace TOKEN)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer TOKEN"
```

## 📝 Request/Response Examples

### Submit Achievement
```http
POST /api/student/achievements
Authorization: Bearer <token>
Content-Type: application/json

{
  "achievement_type": "competition",
  "title": "Programming Competition Winner",
  "description": "Won first place in national programming competition",
  "details": {
    "competition_name": "National Programming Contest",
    "competition_level": "national",
    "rank": 1,
    "participants": 100
  },
  "attachments": [
    {
      "file_name": "certificate.pdf",
      "file_url": "/uploads/certificate.pdf",
      "file_type": "application/pdf"
    }
  ],
  "tags": ["programming", "competition", "national"],
  "points": 100.0
}
```

### Verify Achievement
```http
POST /api/lecturer/achievements/{id}/verify
Authorization: Bearer <lecturer_token>
Content-Type: application/json

{
  "approved": true
}
```

### Reject Achievement
```http
POST /api/lecturer/achievements/{id}/verify
Authorization: Bearer <lecturer_token>
Content-Type: application/json

{
  "approved": false,
  "rejection_note": "Documentation incomplete. Please provide proper certificates."
}
```

## 🔧 Development Tools

### Generate Swagger Docs
```bash
# If using swaggo/swag
swag init -g main.go -o ./docs
```

### Validate OpenAPI
```bash
# Using swagger-codegen
swagger-codegen validate -i docs/swagger.yaml
```

### Export Postman Collection
1. Open Postman
2. Import collection
3. Export as Collection v2.1
4. Save to `docs/` folder

## 🚨 Error Codes

| Code | Description | Example |
|------|-------------|---------|
| 200 | Success | Request completed successfully |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid input data |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error |

## 📞 Support

- **Documentation Issues**: Check Swagger UI for latest updates
- **API Issues**: Test with Postman collection first
- **Authentication Issues**: Verify token format and expiry
- **Permission Issues**: Check user role and required permissions

## 🔄 Changelog

### v1.0.0 (Current)
- ✅ Complete API documentation
- ✅ Swagger UI integration
- ✅ Postman collection
- ✅ All FR-001 to FR-011 implemented
- ✅ RBAC system
- ✅ JWT authentication
- ✅ Comprehensive error handling