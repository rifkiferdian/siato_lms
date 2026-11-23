INSERT IGNORE INTO permissions(name,guard_name) VALUES ('academics.view','web'),('academics.manage','web');
INSERT IGNORE INTO role_has_permissions(permission_id,role_id)
SELECT p.id,r.id FROM permissions p JOIN roles r
WHERE r.name IN ('super_admin','academic_admin') AND p.name IN ('academics.view','academics.manage');
