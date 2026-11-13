INSERT IGNORE INTO permissions(name,guard_name) VALUES
('users.view','web'),('users.manage','web'),('roles.view','web'),('roles.manage','web');

INSERT IGNORE INTO role_has_permissions(permission_id,role_id)
SELECT p.id,r.id FROM permissions p JOIN roles r
WHERE r.name='super_admin' AND p.name IN ('users.view','users.manage','roles.view','roles.manage');
