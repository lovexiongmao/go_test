-- ============================================
-- 运维平台权限管理系统初始数据
-- ============================================

-- 1. 插入权限数据 (permissions)
-- 用户管理权限
INSERT INTO permissions (name, display_name, description, resource, action, status, created_at, updated_at) VALUES
('user:create', '创建用户', '创建新用户的权限', 'user', 'create', 1, NOW(), NOW()),
('user:read', '查看用户', '查看用户信息的权限', 'user', 'read', 1, NOW(), NOW()),
('user:update', '更新用户', '更新用户信息的权限', 'user', 'update', 1, NOW(), NOW()),
('user:delete', '删除用户', '删除用户的权限', 'user', 'delete', 1, NOW(), NOW()),

-- 角色管理权限
('role:create', '创建角色', '创建新角色的权限', 'role', 'create', 1, NOW(), NOW()),
('role:read', '查看角色', '查看角色信息的权限', 'role', 'read', 1, NOW(), NOW()),
('role:update', '更新角色', '更新角色信息的权限', 'role', 'update', 1, NOW(), NOW()),
('role:delete', '删除角色', '删除角色的权限', 'role', 'delete', 1, NOW(), NOW()),

-- 权限管理权限
('permission:create', '创建权限', '创建新权限的权限', 'permission', 'create', 1, NOW(), NOW()),
('permission:read', '查看权限', '查看权限信息的权限', 'permission', 'read', 1, NOW(), NOW()),
('permission:update', '更新权限', '更新权限信息的权限', 'permission', 'update', 1, NOW(), NOW()),
('permission:delete', '删除权限', '删除权限的权限', 'permission', 'delete', 1, NOW(), NOW()),

-- 服务器管理权限
('server:create', '创建服务器', '添加新服务器的权限', 'server', 'create', 1, NOW(), NOW()),
('server:read', '查看服务器', '查看服务器信息的权限', 'server', 'read', 1, NOW(), NOW()),
('server:update', '更新服务器', '更新服务器配置的权限', 'server', 'update', 1, NOW(), NOW()),
('server:delete', '删除服务器', '删除服务器的权限', 'server', 'delete', 1, NOW(), NOW()),
('server:execute', '执行命令', '在服务器上执行命令的权限', 'server', 'execute', 1, NOW(), NOW()),

-- 应用部署权限
('deploy:create', '创建部署', '创建新部署任务的权限', 'deploy', 'create', 1, NOW(), NOW()),
('deploy:read', '查看部署', '查看部署信息的权限', 'deploy', 'read', 1, NOW(), NOW()),
('deploy:update', '更新部署', '更新部署配置的权限', 'deploy', 'update', 1, NOW(), NOW()),
('deploy:delete', '删除部署', '删除部署任务的权限', 'deploy', 'delete', 1, NOW(), NOW()),
('deploy:execute', '执行部署', '执行部署任务的权限', 'deploy', 'execute', 1, NOW(), NOW()),

-- 监控查看权限
('monitor:read', '查看监控', '查看系统监控数据的权限', 'monitor', 'read', 1, NOW(), NOW()),
('monitor:alert', '管理告警', '管理监控告警规则的权限', 'monitor', 'alert', 1, NOW(), NOW()),

-- 日志查看权限
('log:read', '查看日志', '查看系统日志的权限', 'log', 'read', 1, NOW(), NOW()),
('log:download', '下载日志', '下载日志文件的权限', 'log', 'download', 1, NOW(), NOW()),

-- 配置管理权限
('config:create', '创建配置', '创建配置项的权限', 'config', 'create', 1, NOW(), NOW()),
('config:read', '查看配置', '查看配置信息的权限', 'config', 'read', 1, NOW(), NOW()),
('config:update', '更新配置', '更新配置项的权限', 'config', 'update', 1, NOW(), NOW()),
('config:delete', '删除配置', '删除配置项的权限', 'config', 'delete', 1, NOW(), NOW()),

-- 审计日志权限
('audit:read', '查看审计', '查看审计日志的权限', 'audit', 'read', 1, NOW(), NOW()),

-- 数据库管理权限
('database:create', '创建数据库', '创建数据库的权限', 'database', 'create', 1, NOW(), NOW()),
('database:read', '查看数据库', '查看数据库信息的权限', 'database', 'read', 1, NOW(), NOW()),
('database:update', '更新数据库', '更新数据库配置的权限', 'database', 'update', 1, NOW(), NOW()),
('database:delete', '删除数据库', '删除数据库的权限', 'database', 'delete', 1, NOW(), NOW()),
('database:backup', '备份数据库', '备份数据库的权限', 'database', 'backup', 1, NOW(), NOW()),
('database:restore', '恢复数据库', '恢复数据库的权限', 'database', 'restore', 1, NOW(), NOW()),

-- 容器管理权限
('container:create', '创建容器', '创建容器的权限', 'container', 'create', 1, NOW(), NOW()),
('container:read', '查看容器', '查看容器信息的权限', 'container', 'read', 1, NOW(), NOW()),
('container:update', '更新容器', '更新容器配置的权限', 'container', 'update', 1, NOW(), NOW()),
('container:delete', '删除容器', '删除容器的权限', 'container', 'delete', 1, NOW(), NOW()),
('container:start', '启动容器', '启动容器的权限', 'container', 'start', 1, NOW(), NOW()),
('container:stop', '停止容器', '停止容器的权限', 'container', 'stop', 1, NOW(), NOW()),
('container:restart', '重启容器', '重启容器的权限', 'container', 'restart', 1, NOW(), NOW());

-- 2. 插入角色数据 (roles)
INSERT INTO roles (name, display_name, description, status, created_at, updated_at) VALUES
('super_admin', '超级管理员', '拥有所有权限的超级管理员角色', 1, NOW(), NOW()),
('admin', '管理员', '拥有大部分管理权限的管理员角色', 1, NOW(), NOW()),
('ops_engineer', '运维工程师', '负责服务器和应用管理的运维工程师角色', 1, NOW(), NOW()),
('developer', '开发人员', '负责应用部署和查看的开发人员角色', 1, NOW(), NOW()),
('viewer', '只读用户', '只能查看信息的只读用户角色', 1, NOW(), NOW()),
('dba', '数据库管理员', '负责数据库管理的DBA角色', 1, NOW(), NOW()),
('sre', 'SRE工程师', '负责系统可靠性和监控的SRE工程师角色', 1, NOW(), NOW());

-- 3. 插入用户数据 (users)
-- 注意：密码字段需要使用bcrypt加密后的值，这里使用示例密码 "123456" 的bcrypt hash
-- 实际使用时应该使用真实的bcrypt加密密码
-- 示例: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy (对应密码: 123456)
INSERT INTO users (name, email, password, status, created_at, updated_at) VALUES
('超级管理员', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('系统管理员', 'manager@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('运维工程师-张三', 'ops1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('运维工程师-李四', 'ops2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('开发人员-王五', 'dev1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('开发人员-赵六', 'dev2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('只读用户-测试', 'viewer@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('DBA-数据库管理员', 'dba@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW()),
('SRE工程师', 'sre@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1, NOW(), NOW());

-- 4. 关联角色和权限 (role_permissions)
-- 超级管理员：拥有所有权限
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'super_admin') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions;

-- 管理员：拥有大部分管理权限（除了删除用户、删除角色、删除权限等危险操作）
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions
WHERE name NOT IN ('user:delete', 'role:delete', 'permission:delete', 'server:delete', 'database:delete');

-- 运维工程师：服务器管理、应用部署、监控、日志、配置管理
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'ops_engineer') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions
WHERE resource IN ('server', 'deploy', 'monitor', 'log', 'config', 'container')
   OR (resource = 'user' AND action = 'read')
   OR (resource = 'role' AND action = 'read')
   OR (resource = 'permission' AND action = 'read');

-- 开发人员：应用部署、查看监控、查看日志、查看配置
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'developer') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions
WHERE (resource = 'deploy' AND action IN ('create', 'read', 'update', 'execute'))
   OR (resource = 'monitor' AND action = 'read')
   OR (resource = 'log' AND action IN ('read', 'download'))
   OR (resource = 'config' AND action IN ('read', 'update'))
   OR (resource = 'container' AND action IN ('read', 'start', 'stop', 'restart'))
   OR (resource = 'user' AND action = 'read');

-- 只读用户：只能查看
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'viewer') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions
WHERE action = 'read';

-- DBA：数据库相关权限 + 查看权限
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'dba') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions
WHERE resource = 'database'
   OR (resource IN ('user', 'server', 'monitor', 'log') AND action = 'read');

-- SRE工程师：监控、日志、服务器查看、容器管理
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT 
    (SELECT id FROM roles WHERE name = 'sre') as role_id,
    id as permission_id,
    NOW(),
    NOW()
FROM permissions
WHERE resource IN ('monitor', 'log', 'container')
   OR (resource = 'server' AND action IN ('read', 'execute'))
   OR (resource = 'config' AND action IN ('read', 'update'))
   OR (resource = 'deploy' AND action IN ('read', 'execute'))
   OR (resource IN ('user', 'role', 'permission') AND action = 'read');

-- 5. 关联用户和角色 (user_roles)
-- 超级管理员用户 -> 超级管理员角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (
    (SELECT id FROM users WHERE email = 'admin@example.com'),
    (SELECT id FROM roles WHERE name = 'super_admin'),
    NOW(),
    NOW()
);

-- 系统管理员用户 -> 管理员角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (
    (SELECT id FROM users WHERE email = 'manager@example.com'),
    (SELECT id FROM roles WHERE name = 'admin'),
    NOW(),
    NOW()
);

-- 运维工程师用户 -> 运维工程师角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES 
(
    (SELECT id FROM users WHERE email = 'ops1@example.com'),
    (SELECT id FROM roles WHERE name = 'ops_engineer'),
    NOW(),
    NOW()
),
(
    (SELECT id FROM users WHERE email = 'ops2@example.com'),
    (SELECT id FROM roles WHERE name = 'ops_engineer'),
    NOW(),
    NOW()
);

-- 开发人员用户 -> 开发人员角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES 
(
    (SELECT id FROM users WHERE email = 'dev1@example.com'),
    (SELECT id FROM roles WHERE name = 'developer'),
    NOW(),
    NOW()
),
(
    (SELECT id FROM users WHERE email = 'dev2@example.com'),
    (SELECT id FROM roles WHERE name = 'developer'),
    NOW(),
    NOW()
);

-- 只读用户 -> 只读用户角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (
    (SELECT id FROM users WHERE email = 'viewer@example.com'),
    (SELECT id FROM roles WHERE name = 'viewer'),
    NOW(),
    NOW()
);

-- DBA用户 -> DBA角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (
    (SELECT id FROM users WHERE email = 'dba@example.com'),
    (SELECT id FROM roles WHERE name = 'dba'),
    NOW(),
    NOW()
);

-- SRE工程师用户 -> SRE工程师角色
INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
VALUES (
    (SELECT id FROM users WHERE email = 'sre@example.com'),
    (SELECT id FROM roles WHERE name = 'sre'),
    NOW(),
    NOW()
);

-- ============================================
-- 数据说明
-- ============================================
-- 1. 权限表包含47个权限，涵盖：
--    - 用户管理（4个）
--    - 角色管理（4个）
--    - 权限管理（4个）
--    - 服务器管理（5个）
--    - 应用部署（5个）
--    - 监控管理（2个）
--    - 日志管理（2个）
--    - 配置管理（4个）
--    - 审计日志（1个）
--    - 数据库管理（6个）
--    - 容器管理（7个）
--
-- 2. 角色表包含7个角色：
--    - super_admin: 超级管理员（所有权限）
--    - admin: 管理员（大部分权限）
--    - ops_engineer: 运维工程师（服务器、部署、监控等）
--    - developer: 开发人员（部署、查看等）
--    - viewer: 只读用户（仅查看权限）
--    - dba: 数据库管理员（数据库相关权限）
--    - sre: SRE工程师（监控、日志、可靠性相关）
--
-- 3. 用户表包含9个示例用户，密码均为 "123456"（bcrypt加密）
--    实际使用时请修改密码！
--
-- 4. 关联关系：
--    - 每个用户都关联了对应的角色
--    - 每个角色都关联了相应的权限
--
-- 使用说明：
-- 1. 执行此SQL前，请确保表结构已创建（通过GORM AutoMigrate或手动创建）
-- 2. 密码字段使用的是示例bcrypt hash，实际使用时请使用真实的密码加密
-- 3. 可以根据实际需求调整权限、角色和用户的配置
