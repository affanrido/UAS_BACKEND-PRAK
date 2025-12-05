-- Seed data untuk testing

-- Insert roles
INSERT INTO roles (id, name, description) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'admin', 'Administrator dengan akses penuh'),
    ('550e8400-e29b-41d4-a716-446655440002', 'lecturer', 'Dosen yang dapat mengelola mahasiswa'),
    ('550e8400-e29b-41d4-a716-446655440003', 'student', 'Mahasiswa dengan akses terbatas')
ON CONFLICT (name) DO NOTHING;

-- Insert permissions
INSERT INTO permissions (id, name, resource, action, description) VALUES
    ('660e8400-e29b-41d4-a716-446655440001', 'user.read', 'user', 'read', 'Membaca data user'),
    ('660e8400-e29b-41d4-a716-446655440002', 'user.write', 'user', 'write', 'Menulis data user'),
    ('660e8400-e29b-41d4-a716-446655440003', 'user.delete', 'user', 'delete', 'Menghapus data user'),
    ('660e8400-e29b-41d4-a716-446655440004', 'student.read', 'student', 'read', 'Membaca data mahasiswa'),
    ('660e8400-e29b-41d4-a716-446655440005', 'student.write', 'student', 'write', 'Menulis data mahasiswa'),
    ('660e8400-e29b-41d4-a716-446655440006', 'lecturer.read', 'lecturer', 'read', 'Membaca data dosen'),
    ('660e8400-e29b-41d4-a716-446655440007', 'lecturer.write', 'lecturer', 'write', 'Menulis data dosen'),
    ('660e8400-e29b-41d4-a716-446655440008', 'achievement.read', 'achievement', 'read', 'Membaca data prestasi'),
    ('660e8400-e29b-41d4-a716-446655440009', 'achievement.write', 'achievement', 'write', 'Menulis data prestasi'),
    ('660e8400-e29b-41d4-a716-446655440010', 'achievement.verify', 'achievement', 'verify', 'Verifikasi prestasi')
ON CONFLICT (name) DO NOTHING;

-- Assign permissions to admin role (all permissions)
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440002'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440003'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440004'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440005'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440006'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440007'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440008'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440009'),
    ('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440010')
ON CONFLICT DO NOTHING;

-- Assign permissions to lecturer role
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440004'),
    ('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440008'),
    ('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440010')
ON CONFLICT DO NOTHING;

-- Assign permissions to student role
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('550e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440008'),
    ('550e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440009')
ON CONFLICT DO NOTHING;

-- Insert test users (password: "password123" untuk semua)
INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active) VALUES
    (
        '770e8400-e29b-41d4-a716-446655440001',
        'admin',
        'admin@example.com',
        '$2a$10$xhkEhXN/7qZlAF.H368dXeO9sNgJzBpGn3MUDlUB1JqZxhCeo2v9W',
        'Administrator',
        '550e8400-e29b-41d4-a716-446655440001',
        true
    ),
    (
        '770e8400-e29b-41d4-a716-446655440002',
        'lecturer1',
        'lecturer@example.com',
        '$2a$10$xhkEhXN/7qZlAF.H368dXeO9sNgJzBpGn3MUDlUB1JqZxhCeo2v9W',
        'Dr. John Doe',
        '550e8400-e29b-41d4-a716-446655440002',
        true
    ),
    (
        '770e8400-e29b-41d4-a716-446655440003',
        'student1',
        'student@example.com',
        '$2a$10$xhkEhXN/7qZlAF.H368dXeO9sNgJzBpGn3MUDlUB1JqZxhCeo2v9W',
        'Jane Smith',
        '550e8400-e29b-41d4-a716-446655440003',
        true
    )
ON CONFLICT (username) DO NOTHING;

-- Insert lecturers
INSERT INTO lecturers (id, user_id, lecturer_id, department) VALUES
    (
        '880e8400-e29b-41d4-a716-446655440001',
        '770e8400-e29b-41d4-a716-446655440002',
        'LEC001',
        'Computer Science'
    )
ON CONFLICT (lecturer_id) DO NOTHING;

-- Insert students
INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id) VALUES
    (
        '990e8400-e29b-41d4-a716-446655440001',
        '770e8400-e29b-41d4-a716-446655440003',
        'STD001',
        'Teknik Informatika',
        '2021',
        '880e8400-e29b-41d4-a716-446655440001'
    )
ON CONFLICT (student_id) DO NOTHING;
