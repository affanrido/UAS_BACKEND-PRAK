# UAS Backend Deployment Guide

Comprehensive guide untuk deployment UAS Backend API dengan dokumentasi Swagger.

## 🚀 Quick Deployment

### 1. Prerequisites
```bash
# Go 1.21+
go version

# PostgreSQL 12+
psql --version

# MongoDB 4.4+
mongod --version
```

### 2. Clone & Setup
```bash
# Clone repository
git clone <repository-url>
cd UAS_BACKEND

# Install dependencies
go mod tidy

# Setup environment
cp .env.example .env
# Edit .env with your database credentials
```

### 3. Database Setup
```bash
# PostgreSQL
createdb uas_backend
psql -U postgres -d uas_backend -f database/schema.sql
psql -U postgres -d uas_backend -f database/seed.sql

# MongoDB (auto-created on first connection)
```

### 4. Run Application
```bash
# Development
go run main.go

# Production build
go build -o uas-backend .
./uas-backend
```

### 5. Access Documentation
- **API Docs**: http://localhost:8080/swagger/
- **Health Check**: http://localhost:8080/health
- **Test Login**: Use Postman collection in `docs/`

## 📋 Environment Configuration

### Required Environment Variables
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=uas_backend

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=uas_backend

# JWT Configuration
JWT_SECRET=your-super-secret-key-here
JWT_EXPIRY=24h

# Server Configuration
PORT=8080
ENV=development

# File Upload Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760  # 10MB in bytes
```

### Optional Environment Variables
```env
# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# CORS Configuration
CORS_ORIGINS=*
CORS_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_HEADERS=*

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Cache Configuration
CACHE_TTL=5m
CACHE_CLEANUP_INTERVAL=10m
```

## 🐳 Docker Deployment

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o uas-backend .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/uas-backend .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/database ./database

EXPOSE 8080
CMD ["./uas-backend"]
```

### Docker Compose
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - MONGO_URI=mongodb://mongo:27017
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
      - mongo
    volumes:
      - ./uploads:/root/uploads

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: uas_backend
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  mongo:
    image: mongo:6-jammy
    environment:
      MONGO_INITDB_DATABASE: uas_backend
    volumes:
      - mongo_data:/data/db
    ports:
      - "27017:27017"

volumes:
  postgres_data:
  mongo_data:
```

### Deploy with Docker
```bash
# Build and run
docker-compose up -d

# Check logs
docker-compose logs -f app

# Access documentation
open http://localhost:8080/swagger/
```

## ☁️ Cloud Deployment

### Heroku Deployment
```bash
# Install Heroku CLI
# Create Heroku app
heroku create uas-backend-api

# Add PostgreSQL addon
heroku addons:create heroku-postgresql:mini

# Add MongoDB addon
heroku addons:create mongolab:sandbox

# Set environment variables
heroku config:set JWT_SECRET=your-secret-key
heroku config:set ENV=production

# Deploy
git push heroku main

# Run database migrations
heroku run psql $DATABASE_URL -f database/schema.sql
heroku run psql $DATABASE_URL -f database/seed.sql
```

### AWS Deployment (ECS)
```yaml
# task-definition.json
{
  "family": "uas-backend",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "uas-backend",
      "image": "your-account.dkr.ecr.region.amazonaws.com/uas-backend:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "DB_HOST",
          "value": "your-rds-endpoint"
        },
        {
          "name": "MONGO_URI",
          "value": "mongodb://your-documentdb-endpoint:27017"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/uas-backend",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

## 🔧 Production Configuration

### Nginx Reverse Proxy
```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Swagger documentation
    location /swagger/ {
        proxy_pass http://localhost:8080/swagger/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Static files
    location /docs/ {
        proxy_pass http://localhost:8080/docs/;
        expires 1d;
        add_header Cache-Control "public, immutable";
    }
}
```

### SSL Configuration (Let's Encrypt)
```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d api.yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Systemd Service
```ini
# /etc/systemd/system/uas-backend.service
[Unit]
Description=UAS Backend API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/uas-backend
ExecStart=/opt/uas-backend/uas-backend
Restart=always
RestartSec=5
Environment=ENV=production

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl enable uas-backend
sudo systemctl start uas-backend
sudo systemctl status uas-backend
```

## 📊 Monitoring & Logging

### Health Check Endpoint
```bash
# Basic health check
curl http://localhost:8080/health

# Response
{"status":"ok"}
```

### Application Metrics
```go
// Add to main.go for monitoring
import (
    "github.com/gofiber/fiber/v2/middleware/monitor"
)

// Add monitoring endpoint
app.Get("/metrics", monitor.New())
```

### Log Configuration
```go
// Enhanced logging
app.Use(logger.New(logger.Config{
    Format: "${time} ${status} - ${method} ${path} - ${latency}\n",
    TimeFormat: "2006-01-02 15:04:05",
    Output: os.Stdout,
}))
```

## 🔒 Security Considerations

### Production Security Checklist
- [ ] Use strong JWT secret (32+ characters)
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS properly
- [ ] Set up rate limiting
- [ ] Use environment variables for secrets
- [ ] Enable database connection encryption
- [ ] Set up firewall rules
- [ ] Regular security updates
- [ ] Monitor for vulnerabilities
- [ ] Backup strategy

### Environment Security
```env
# Use strong secrets
JWT_SECRET=your-very-long-and-random-secret-key-here-32-chars-minimum

# Restrict CORS in production
CORS_ORIGINS=https://yourdomain.com,https://app.yourdomain.com

# Database security
DB_SSL_MODE=require
MONGO_SSL=true
```

## 📈 Performance Optimization

### Database Optimization
```sql
-- Add indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_achievement_references_student_id ON achievement_references(student_id);
CREATE INDEX idx_achievement_references_status ON achievement_references(status);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
```

### Application Optimization
```go
// Connection pooling
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)

// Enable compression
app.Use(compress.New())

// Cache static files
app.Static("/docs", "./docs", fiber.Static{
    Compress: true,
    MaxAge:   86400, // 1 day
})
```

## 🧪 Testing in Production

### Smoke Tests
```bash
#!/bin/bash
# smoke-test.sh

BASE_URL="https://api.yourdomain.com"

echo "Running smoke tests..."

# Health check
curl -f "$BASE_URL/health" || exit 1

# Swagger documentation
curl -f "$BASE_URL/swagger/" || exit 1

# Login test
TOKEN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"identifier":"admin@example.com","password":"password123"}' \
  | jq -r '.data.token')

if [ "$TOKEN" != "null" ]; then
  echo "✅ All tests passed"
else
  echo "❌ Login test failed"
  exit 1
fi
```

### Load Testing
```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:8080/health

# Using wrk
wrk -t12 -c400 -d30s http://localhost:8080/health
```

## 📞 Troubleshooting

### Common Issues

#### 1. Database Connection Failed
```bash
# Check PostgreSQL
sudo systemctl status postgresql
sudo -u postgres psql -c "SELECT version();"

# Check MongoDB
sudo systemctl status mongod
mongo --eval "db.adminCommand('ismaster')"
```

#### 2. Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080
sudo kill -9 <PID>
```

#### 3. Swagger UI Not Loading
```bash
# Check if docs directory exists
ls -la docs/

# Check file permissions
chmod -R 755 docs/

# Verify Swagger files
curl http://localhost:8080/swagger/swagger.yaml
```

#### 4. JWT Token Issues
```bash
# Check JWT secret length
echo $JWT_SECRET | wc -c  # Should be 32+ characters

# Verify token format
echo "TOKEN" | base64 -d
```

### Debug Mode
```bash
# Run with debug logging
ENV=development LOG_LEVEL=debug go run main.go

# Enable Fiber debug mode
FIBER_DEBUG=true go run main.go
```

## 📚 Additional Resources

- **API Documentation**: http://localhost:8080/swagger/
- **Postman Collection**: `docs/UAS_Backend_API.postman_collection.json`
- **Database Schema**: `database/schema.sql`
- **Test Data**: `database/seed.sql`
- **Environment Template**: `.env.example`

## 🆘 Support

For deployment issues:
1. Check logs: `docker-compose logs -f app`
2. Verify environment variables
3. Test database connections
4. Check firewall/security groups
5. Validate SSL certificates