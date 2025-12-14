# Database Documentation

## Mengapa Ada File .sql? 🗄️

File .sql dalam folder `database/` adalah **script database** yang diperlukan untuk setup dan inisialisasi sistem UAS Backend. Berikut penjelasan lengkapnya:

## 📁 File Structure

```
database/
├── schema.sql          # 🏗️ Database schema (struktur tabel)
└── seed.sql           # 🌱 Data awal untuk testing
```

## 🏗️ schema.sql - Database Schema

**Tujuan**: Membuat struktur database PostgreSQL yang diperlukan sistem

**Isi**:
- **Tabel Users** - Data pengguna (admin, lecturer, student)
- **Tabel Roles** - Role sistem (admin, lecturer, student)
- **Tabel Permissions** - Hak akses sistem
- **Tabel Role_Permissions** - Mapping role ke permissions
- **Tabel Students** - Data mahasiswa
- **Tabel Lecturers** - Data dosen
- **Tabel Achievement_References** - Referensi prestasi (PostgreSQL)
- **Tabel Notifications** - Sistem notifikasi

**Contoh**:
```sql
-- Tabel users untuk authentication
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🌱 seed.sql - Data Awal

**Tujuan**: Mengisi database dengan data testing yang diperlukan

**Isi**:
- **3 Roles**: admin, lecturer, student
- **Permissions**: Semua hak akses sistem (user.read, achievement.write, dll)
- **Role-Permission Mapping**: Menentukan role mana yang punya akses apa
- **Test Users**: 3 user untuk testing (admin, lecturer, student)
- **Sample Data**: Data mahasiswa dan dosen untuk testing

**Test Users**:
```sql
-- User admin untuk testing
INSERT INTO users (id, username, email, password_hash, full_name, role_id) VALUES
('770e8400-e29b-41d4-a716-446655440001', 'admin', 'admin@example.com', 
 '$2a$10$...', 'Administrator', '550e8400-e29b-41d4-a716-446655440001');

-- User lecturer untuk testing  
INSERT INTO users (id, username, email, password_hash, full_name, role_id) VALUES
('770e8400-e29b-41d4-a716-446655440002', 'lecturer1', 'lecturer@example.com',
 '$2a$10$...', 'Dr. John Lecturer', '550e8400-e29b-41d4-a716-446655440002');

-- User student untuk testing
INSERT INTO users (id, username, email, password_hash, full_name, role_id) VALUES
('770e8400-e29b-41d4-a716-446655440003', 'student1', 'student@example.com',
 '$2a$10$...', 'Jane Student', '550e8400-e29b-41d4-a716-446655440003');
```

## 🚀 Cara Menggunakan

### 1. Setup Database Baru
```bash
# Buat database PostgreSQL
createdb uas_backend

# Jalankan schema untuk membuat tabel
psql -U postgres -d uas_backend -f database/schema.sql

# Jalankan seed untuk mengisi data awal
psql -U postgres -d uas_backend -f database/seed.sql
```

### 2. Reset Database
```bash
# Drop dan buat ulang database
dropdb uas_backend
createdb uas_backend

# Jalankan ulang schema dan seed
psql -U postgres -d uas_backend -f database/schema.sql
psql -U postgres -d uas_backend -f database/seed.sql
```

### 3. Update Schema
```bash
# Hanya jalankan schema jika ada perubahan struktur
psql -U postgres -d uas_backend -f database/schema.sql
```

## 🔑 Test Credentials

Setelah menjalankan seed.sql, tersedia test users:

| Role | Username | Email | Password |
|------|----------|-------|----------|
| Admin | admin | admin@example.com | password123 |
| Lecturer | lecturer1 | lecturer@example.com | password123 |
| Student | student1 | student@example.com | password123 |

## 🏛️ Database Architecture

### PostgreSQL (Relational Data)
- **Users & Authentication** - Data user, roles, permissions
- **Student & Lecturer Profiles** - Data akademik
- **Achievement References** - Referensi prestasi dengan status
- **Notifications** - Sistem notifikasi

### MongoDB (Document Data)  
- **Achievements** - Detail prestasi dengan struktur fleksibel
- **File Attachments** - Metadata file upload
- **Dynamic Fields** - Data yang berubah-ubah per tipe prestasi

## 🔧 Database Tools

### Makefile Commands
```bash
# Setup database
make db-setup

# Reset database  
make db-reset

# Backup database
make db-backup

# Restore database
make db-restore
```

### Manual Commands
```bash
# Connect to database
psql -U postgres -d uas_backend

# Check tables
\dt

# Check data
SELECT * FROM users;
SELECT * FROM roles;
SELECT * FROM permissions;
```

## 🛡️ Security Notes

1. **Password Hashing**: Semua password di-hash dengan bcrypt
2. **UUID Primary Keys**: Menggunakan UUID untuk security
3. **Soft Delete**: Data tidak dihapus permanen
4. **Audit Trail**: Timestamp created_at dan updated_at
5. **Role-Based Access**: Sistem permission yang ketat

## 📊 Database Monitoring

### Performance Queries
```sql
-- Check table sizes
SELECT schemaname,tablename,attname,n_distinct,correlation 
FROM pg_stats WHERE tablename = 'users';

-- Check active connections
SELECT * FROM pg_stat_activity WHERE datname = 'uas_backend';

-- Check slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 10;
```

### Backup Strategy
```bash
# Daily backup
pg_dump -U postgres uas_backend > backup_$(date +%Y%m%d).sql

# Compressed backup
pg_dump -U postgres -Fc uas_backend > backup_$(date +%Y%m%d).dump
```

## 🔄 Migration Strategy

Untuk perubahan schema di production:

1. **Backup** database terlebih dahulu
2. **Test** migration di development
3. **Apply** dengan downtime minimal
4. **Verify** data integrity
5. **Rollback** plan jika ada masalah

File .sql ini adalah **fondasi sistem** yang memastikan database siap digunakan dengan data testing yang lengkap! 🎯