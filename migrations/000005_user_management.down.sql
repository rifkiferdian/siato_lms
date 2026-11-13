DELETE rp FROM role_has_permissions rp JOIN permissions p ON p.id=rp.permission_id
WHERE p.name IN ('users.view','users.manage','roles.view','roles.manage');
DELETE FROM permissions WHERE name IN ('users.view','users.manage','roles.view','roles.manage');
