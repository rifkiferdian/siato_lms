DELETE FROM role_has_permissions WHERE permission_id IN (SELECT id FROM permissions WHERE name IN ('students.view','students.manage','lecturers.view','lecturers.manage'));
DELETE FROM permissions WHERE name IN ('students.view','students.manage','lecturers.view','lecturers.manage');
DROP INDEX idx_lecturers_status_source ON lecturers;
DROP INDEX idx_students_status_source ON students;
ALTER TABLE lecturers DROP COLUMN last_synced_at, DROP COLUMN data_source, DROP COLUMN academic_status;
ALTER TABLE students DROP COLUMN last_synced_at, DROP COLUMN data_source, DROP COLUMN academic_status;
