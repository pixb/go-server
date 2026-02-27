-- 系统设置数据
INSERT OR IGNORE INTO system_setting (name, value, description) VALUES
('instance_name', 'Echo Demo', '实例名称'),
('instance_version', '0.1.0', '实例版本'),
('schema_version', '0.1.0', '数据库 schema 版本'),
('admin_email', 'admin@example.com', '管理员邮箱'),
('site_url', 'http://localhost:8081', '站点 URL');

-- 用户数据（密码都是 password123）
INSERT OR IGNORE INTO users (id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at) VALUES
(1, 'admin', 'Admin User', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138000', 'admin@example.com', 'admin', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(2, 'user1', 'Test User 1', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138001', 'user1@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(3, 'user2', 'Test User 2', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138002', 'user2@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(4, 'user3', 'Test User 3', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138003', 'user3@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(5, 'user4', 'Test User 4', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138004', 'user4@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(6, 'user5', 'Test User 5', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138005', 'user5@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(7, 'user6', 'Test User 6', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138006', 'user6@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(8, 'user7', 'Test User 7', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138007', 'user7@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(9, 'user8', 'Test User 8', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138008', 'user8@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(10, 'user9', 'Test User 9', '$2a$10$z6T0xV7yqJ5B7K8uG9x1Ou6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q6Q', '13800138009', 'user9@example.com', 'user', '2027-02-27 00:00:00', '2026-02-27 00:00:00', '2026-02-27 00:00:00');

-- 刷新令牌数据
INSERT OR IGNORE INTO refresh_tokens (user_id, token, expires_at, revoked, created_at, updated_at) VALUES
(1, 'admin-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(2, 'user1-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(3, 'user2-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(4, 'user3-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(5, 'user4-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(6, 'user5-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(7, 'user6-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(8, 'user7-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(9, 'user8-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00'),
(10, 'user9-refresh-token-1', '2027-02-27 00:00:00', 0, '2026-02-27 00:00:00', '2026-02-27 00:00:00');
