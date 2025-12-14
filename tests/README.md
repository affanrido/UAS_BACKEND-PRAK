# Unit Testing Documentation

This directory contains comprehensive unit tests for the UAS Backend project following the requirements:

## 6.1 Unit Testing Requirements
- ✅ Test individual functions and methods
- ✅ Mock external dependencies
- ✅ Cover success cases, error cases, and edge cases
- ✅ Use testify framework for assertions and mocks

## Test Structure

```
tests/
├── mocks/                          # Mock implementations
│   ├── auth_repository_mock.go     # Auth repository mock
│   ├── auth_service_mock.go        # Auth service mock
│   ├── achievement_repository_mock.go # Achievement repository mock
│   ├── user_repository_mock.go     # User repository mock
│   ├── statistics_repository_mock.go # Statistics repository mock
│   ├── notification_repository_mock.go # Notification repository mock
│   └── notification_service_mock.go # Notification service mock
├── service/                        # Service layer tests
│   ├── auth_service_test.go        # Authentication service tests
│   ├── achievement_service_test.go # Achievement service tests
│   ├── user_service_test.go        # User management service tests
│   ├── statistics_service_test.go  # Statistics service tests
│   └── notification_service_test.go # Notification service tests
├── middleware/                     # Middleware tests
│   ├── auth_middleware_test.go     # Authentication middleware tests
│   └── rbac_middleware_test.go     # RBAC middleware tests
├── integration/                    # Integration tests
│   └── auth_handler_test.go        # HTTP handler integration tests
└── README.md                       # This file
```

## Test Coverage

### Service Layer Tests
1. **AuthService** (`auth_service_test.go`)
   - Login with valid credentials
   - Login with invalid credentials
   - Login with inactive user
   - Token validation (valid/invalid/expired)
   - Empty credentials validation

2. **AchievementService** (`achievement_service_test.go`)
   - Submit achievement (success/validation errors)
   - Submit for verification (success/wrong status)
   - Delete achievement (success/wrong status)
   - User authorization checks

3. **UserService** (`user_service_test.go`)
   - Create user (success/validation/conflicts)
   - Update user (success/not found/conflicts)
   - Delete user (success/not found)
   - Role assignment
   - Profile management (student/lecturer)

4. **StatisticsService** (`statistics_service_test.go`)
   - Student statistics (own achievements)
   - Lecturer statistics (advisee achievements)
   - Admin statistics (all achievements)
   - Achievement trends with different roles
   - Input validation and error handling

5. **NotificationService** (`notification_service_test.go`)
   - Create notifications (submitted/verified/rejected)
   - Get user notifications with pagination
   - Mark notifications as read
   - Unread count functionality

### Middleware Tests
1. **AuthMiddleware** (`auth_middleware_test.go`)
   - Valid token authentication
   - Missing/invalid token handling
   - Malformed authorization headers

2. **RBACMiddleware** (`rbac_middleware_test.go`)
   - Permission-based access control
   - Role-based access control
   - Missing user context handling

### Integration Tests
1. **AuthHandler** (`auth_handler_test.go`)
   - Login endpoint integration
   - Request validation
   - Response format verification

## Running Tests

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Specific Test Packages
```bash
# Service tests only
go test ./tests/service -v

# Middleware tests only
go test ./tests/middleware -v

# Integration tests only
go test ./tests/integration -v
```

### Run Individual Test Files
```bash
go test ./tests/service/auth_service_test.go -v
go test ./tests/service/user_service_test.go -v
```

### Run with Coverage
```bash
go test ./tests/... -v -cover
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Using Test Runner Script
```bash
go run run_tests.go
```

## Mock Usage

All external dependencies are mocked using testify/mock:

- **Repository Layer**: Database operations are mocked
- **Service Dependencies**: Cross-service calls are mocked
- **External APIs**: Third-party integrations are mocked

Example mock usage:
```go
mockRepo := new(mocks.MockUserRepository)
mockRepo.On("GetUserByID", userID).Return(expectedUser, nil)
userService := service.NewUserService(mockRepo)

result, err := userService.GetUserByID(userID)

assert.NoError(t, err)
assert.Equal(t, expectedUser.ID, result.User.ID)
mockRepo.AssertExpectations(t)
```

## Test Data

Tests use:
- Generated UUIDs for entity IDs
- Realistic test data that matches production models
- Edge cases (empty strings, nil values, boundary conditions)
- Error scenarios (database failures, validation errors)

## Dependencies

Required testing packages (already in go.mod):
- `github.com/stretchr/testify` - Assertions and mocks
- `github.com/DATA-DOG/go-sqlmock` - SQL database mocking
- `github.com/golang/mock` - Mock generation

## Best Practices Followed

1. **Arrange-Act-Assert Pattern**: Clear test structure
2. **Descriptive Test Names**: Self-documenting test purposes
3. **Independent Tests**: No test dependencies
4. **Mock Verification**: Assert all mock expectations
5. **Error Testing**: Comprehensive error case coverage
6. **Edge Case Testing**: Boundary and null value testing

## Continuous Integration

Tests are designed to run in CI/CD pipelines:
- No external dependencies required
- Fast execution with mocked dependencies
- Deterministic results
- Clear pass/fail indicators