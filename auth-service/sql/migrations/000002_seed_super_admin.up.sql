INSERT INTO users (email, password_hash, full_name, phone_number, role)
VALUES ('admin@bus.com', '$2a$10$WRPsPboRFEu/thP0mLC5QOq7zo.HrISnDHFRZyZzYncux/7A43LlO', 'Super Admin', '+251911111111', 'ROLE_SUPER_ADMIN')
ON CONFLICT (email) DO NOTHING;
