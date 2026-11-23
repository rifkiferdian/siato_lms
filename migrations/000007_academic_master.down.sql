DELETE FROM role_has_permissions WHERE permission_id IN (SELECT id FROM permissions WHERE name IN ('academics.view','academics.manage'));
DELETE FROM permissions WHERE name IN ('academics.view','academics.manage');
