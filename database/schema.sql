-- Create tables for authentication and achievement system

-- 3.1.1 Tabel users
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

-- 3.1.2 Tabel roles
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3.1.3 Tabel permissions
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

-- 3.1.4 Tabel role_permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- 3.1.5 Tabel students
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    student_id VARCHAR(20) UNIQUE NOT NULL,
    program_study VARCHAR(100),
    academic_year VARCHAR(10),
    advisor_id UUID REFERENCES lecturers(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3.1.6 Tabel lecturers
CREATE TABLE IF NOT EXISTS lecturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lecturer_id VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3.1.7 Tabel achievement_references
CREATE TABLE IF NOT EXISTS achievement_references (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    mongo_achievement_id VARCHAR(24) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('draft', 'submitted', 'verified', 'rejected')) DEFAULT 'draft',
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id),
    rejection_note TEXT,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add foreign key constraint for users -> roles
ALTER TABLE users ADD CONSTRAINT fk_users_role 
    FOREIGN KEY (role_id) REFERENCES roles(id);

-- 3.1.8 Tabel notifications
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    related_id UUID,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_students_user_id ON students(user_id);
CREATE INDEX IF NOT EXISTS idx_students_student_id ON students(student_id);
CREATE INDEX IF NOT EXISTS idx_students_advisor_id ON students(advisor_id);
CREATE INDEX IF NOT EXISTS idx_lecturers_user_id ON lecturers(user_id);
CREATE INDEX IF NOT EXISTS idx_lecturers_lecturer_id ON lecturers(lecturer_id);
CREATE INDEX IF NOT EXISTS idx_achievement_references_student_id ON achievement_references(student_id);
CREATE INDEX IF NOT EXISTS idx_achievement_references_mongo_id ON achievement_references(mongo_achievement_id);
CREATE INDEX IF NOT EXISTS idx_achievement_references_status ON achievement_references(status);
CREATE INDEX IF NOT EXISTS idx_achievement_references_is_deleted ON achievement_references(is_deleted);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
