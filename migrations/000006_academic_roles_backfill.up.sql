INSERT IGNORE INTO model_has_roles(role_id,user_id)
SELECT r.id,s.user_id FROM students s JOIN roles r ON r.name='student';

INSERT IGNORE INTO model_has_roles(role_id,user_id)
SELECT r.id,l.user_id FROM lecturers l JOIN roles r ON r.name='lecturer';

CREATE TRIGGER trg_students_assign_role
AFTER INSERT ON students FOR EACH ROW
INSERT IGNORE INTO model_has_roles(role_id,user_id)
SELECT id,NEW.user_id FROM roles WHERE name='student';

CREATE TRIGGER trg_lecturers_assign_role
AFTER INSERT ON lecturers FOR EACH ROW
INSERT IGNORE INTO model_has_roles(role_id,user_id)
SELECT id,NEW.user_id FROM roles WHERE name='lecturer';
