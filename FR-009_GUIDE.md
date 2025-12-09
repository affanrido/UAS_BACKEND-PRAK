# FR-009: Manage Users - Implementation Guide

## Overview

Implementasi lengkap fitur untuk admin melakukan CRUD users dan assign roles, termasuk set student/lecturer profile dan set advisor untuk mahasiswa.

---

## Flow Implementation

### 1. ✅ Create/Update/Delete User
**Implementation:** Full CRUD operations

**Create User:**
```http
POST /api/admin/users
```

**Update User:**
```http
PUT /api/admin/users/:id
```

**Delete User:**
```http
DELETE /api/admin/users/:id
```

---

### 2. ✅ Assign Role
**Implementation:** Assign role to user

```http
POST /api/admin/users/:id/assign-role
```

---

### 3. ✅ Set Student/Lecturer Profile
**Implementation:** Set or update student/lecturer profile

**Set Student Profile:**
```http
POST /api/admin/users/:id/student-profile
```

**Set Lecturer Profile:**
```http
POST /api/admin/users/:id/lecturer-profile
```

---

### 4. ✅ Set Advisor untuk Mahasiswa
**Implementation:** Assign advisor to student

```http
POST /api/admin/users/:id/set-advisor
```

---

## API Endpoints

### 1. Create User

**POST /api/admin/users**

Create new user with optional student/lecturer profile.

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "role_id": "550e8400-e29b-41d4-a716-446655440003",
  "is_active": true,
  "student_data": {
    "student_id": "STD001",
    "program_study": "Computer Science",
    "academic_year": "2021",
    "advisor_id": "660e8400-e29b-41d4-a716-446655440001"
  }
}
```

**Success Response (201):**
```json
{
  "message": "User created successfully",
  "data": {
    "user": {
      "id": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role_id": "uuid",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "student": {
      "id": "uuid",
      "user_id": "uuid",
      "student_id": "STD001",
      "program_study": "Computer Science",
      "academic_year": "2021",
      "advisor_id": "uuid",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "role": {
      "id": "uuid",
      "name": "student",
      "description": "Student role"
    }
  }
}
```

---

### 2. Get All Users

**GET /api/admin/users**

Get all users with pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)

**Success Response (200):**
```json
{
  "message": "Users retrieved successfully",
  "data": [
    {
      "user": {...},
      "student": {...},
      "lecturer": null,
      "role": {...}
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 25,
    "total_pages": 3
  }
}
```

---

### 3. Get User by ID

**GET /api/admin/users/:id**

Get user details by ID.

**Success Response (200):**
```json
{
  "message": "User retrieved successfully",
  "data": {
    "user": {...},
    "student": {...},
    "lecturer": null,
    "role": {...}
  }
}
```

---

### 4. Update User

**PUT /api/admin/users/:id**

Update existing user.

**Request Body (all fields optional):**
```json
{
  "username": "new_username",
  "email": "newemail@example.com",
  "password": "newpassword123",
  "full_name": "New Full Name",
  "role_id": "new-role-uuid",
  "is_active": false
}
```

**Success Response (200):**
```json
{
  "message": "User updated successfully",
  "data": {
    "user": {...},
    "student": {...},
    "role": {...}
  }
}
```

---

### 5. Delete User

**DELETE /api/admin/users/:id**

Delete user (cascade deletes student/lecturer profiles).

**Success Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

---

### 6. Assign Role

**POST /api/admin/users/:id/assign-role**

Assign role to user.

**Request Body:**
```json
{
  "role_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

**Success Response (200):**
```json
{
  "message": "Role assigned successfully",
  "data": {
    "user": {...},
    "role": {
      "id": "uuid",
      "name": "lecturer",
      "description": "Lecturer role"
    }
  }
}
```

---

### 7. Set Student Profile

**POST /api/admin/users/:id/student-profile**

Set or update student profile for user.

**Request Body:**
```json
{
  "student_id": "STD001",
  "program_study": "Computer Science",
  "academic_year": "2021",
  "advisor_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

**Success Response (200):**
```json
{
  "message": "Student profile set successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "student_id": "STD001",
    "program_study": "Computer Science",
    "academic_year": "2021",
    "advisor_id": "uuid",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 8. Set Lecturer Profile

**POST /api/admin/users/:id/lecturer-profile**

Set or update lecturer profile for user.

**Request Body:**
```json
{
  "lecturer_id": "LEC001",
  "department": "Computer Science"
}
```

**Success Response (200):**
```json
{
  "message": "Lecturer profile set successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "lecturer_id": "LEC001",
    "department": "Computer Science",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 9. Set Advisor

**POST /api/admin/users/:id/set-advisor**

Set advisor for student.

**Request Body:**
```json
{
  "advisor_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

**Success Response (200):**
```json
{
  "message": "Advisor set successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "student_id": "STD001",
    "program_study": "Computer Science",
    "academic_year": "2021",
    "advisor_id": "660e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Complete Example

### Step 1: Login as Admin

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "admin@example.com",
    "password": "password123"
  }'
```

**Save TOKEN!**

---

### Step 2: Create Student User

```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jane_smith",
    "email": "jane@example.com",
    "password": "password123",
    "full_name": "Jane Smith",
    "role_id": "550e8400-e29b-41d4-a716-446655440003",
    "is_active": true,
    "student_data": {
      "student_id": "STD002",
      "program_study": "Information Systems",
      "academic_year": "2022",
      "advisor_id": "660e8400-e29b-41d4-a716-446655440001"
    }
  }'
```

---

### Step 3: Create Lecturer User

```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "dr_john",
    "email": "drjohn@example.com",
    "password": "password123",
    "full_name": "Dr. John Doe",
    "role_id": "550e8400-e29b-41d4-a716-446655440002",
    "is_active": true,
    "lecturer_data": {
      "lecturer_id": "LEC002",
      "department": "Computer Science"
    }
  }'
```

---

### Step 4: Get All Users

```bash
curl -X GET "http://localhost:8080/api/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer TOKEN"
```

---

### Step 5: Update User

```bash
curl -X PUT http://localhost:8080/api/admin/users/USER_ID \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Jane Doe Smith",
    "is_active": true
  }'
```

---

### Step 6: Assign Role

```bash
curl -X POST http://localhost:8080/api/admin/users/USER_ID/assign-role \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": "550e8400-e29b-41d4-a716-446655440002"
  }'
```

---

### Step 7: Set Advisor

```bash
curl -X POST http://localhost:8080/api/admin/users/STUDENT_USER_ID/set-advisor \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "advisor_id": "LECTURER_ID"
  }'
```

---

### Step 8: Delete User

```bash
curl -X DELETE http://localhost:8080/api/admin/users/USER_ID \
  -H "Authorization: Bearer TOKEN"
```

---

## Validation Rules

### Create User
- ✅ Username required (unique)
- ✅ Email required (unique, valid format)
- ✅ Password required (min 6 characters)
- ✅ Full name required
- ✅ Role ID required (must exist)
- ✅ Student/Lecturer data optional

### Update User
- ✅ All fields optional
- ✅ Username must be unique if changed
- ✅ Email must be unique if changed
- ✅ Password min 6 characters if changed

### Set Student Profile
- ✅ Student ID required
- ✅ Program study required
- ✅ Academic year required
- ✅ Advisor ID required (must be valid lecturer)

### Set Lecturer Profile
- ✅ Lecturer ID required
- ✅ Department required

### Set Advisor
- ✅ Advisor ID required
- ✅ Must be valid lecturer ID
- ✅ User must have student profile

---

## Security Features

1. ✅ **Admin Only** - All endpoints require admin role
2. ✅ **Authentication Required** - JWT token validation
3. ✅ **Password Hashing** - bcrypt for password security
4. ✅ **Unique Constraints** - Username and email must be unique
5. ✅ **Cascade Delete** - Student/lecturer profiles deleted with user
6. ✅ **Input Validation** - All inputs validated

---

## Database Operations

### Users Table
- CREATE: Insert new user
- READ: Get by ID, username, email, or all with pagination
- UPDATE: Update user fields
- DELETE: Delete user

### Students Table
- CREATE: Insert student profile
- READ: Get by user_id
- UPDATE: Update student profile
- DELETE: Cascade on user delete

### Lecturers Table
- CREATE: Insert lecturer profile
- READ: Get by user_id or ID
- UPDATE: Update lecturer profile
- DELETE: Cascade on user delete

---

## Error Handling

| Error | Status | Message |
|-------|--------|---------|
| Not authenticated | 401 | "Unauthorized" |
| Not admin | 403 | "Forbidden" |
| Invalid body | 400 | "Invalid request body" |
| Invalid UUID | 400 | "Invalid user ID" / "Invalid role ID" |
| Username exists | 400 | "username already exists" |
| Email exists | 400 | "email already exists" |
| User not found | 404 | "user not found" |
| Role not found | 400 | "role not found" |
| Student not found | 400 | "student not found" |
| Advisor not found | 400 | "advisor not found" |
| Validation error | 400 | Specific validation message |

---

## Architecture

```
Request (Admin)
  ↓
[RBAC Middleware] ← Check authentication & admin role
  ↓
[Admin Handler] ← Parse request
  ↓
[User Service] ← Business logic & validation
  ↓
[User Repository]
  ├─> [PostgreSQL - users] ← CRUD operations
  ├─> [PostgreSQL - students] ← Student profile
  ├─> [PostgreSQL - lecturers] ← Lecturer profile
  └─> [PostgreSQL - roles] ← Role info
  ↓
Response (User data)
```

---

## Testing Checklist

- [x] Build succeeds
- [x] Create user works
- [x] Update user works
- [x] Delete user works
- [x] Assign role works
- [x] Set student profile works
- [x] Set lecturer profile works
- [x] Set advisor works
- [x] Get all users with pagination
- [x] Get user by ID
- [x] Validation works
- [x] Error handling works
- [x] Admin-only access enforced

---

## Conclusion

FR-009: Manage Users telah **FULLY IMPLEMENTED** dengan:

✅ **All 4 flow steps completed**
✅ **Full CRUD operations**
✅ **Role assignment**
✅ **Student/Lecturer profile management**
✅ **Advisor assignment**
✅ **Complete validation & security**
✅ **Pagination support**
✅ **Error handling**
✅ **Documentation complete**

Sistem siap digunakan untuk admin mengelola users! 🚀
