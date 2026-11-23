DROP INDEX uq_course_curriculum_code ON courses;
CREATE UNIQUE INDEX uq_course_code ON courses(code);
ALTER TABLE courses DROP COLUMN is_active, DROP COLUMN course_type, DROP COLUMN recommended_semester;
ALTER TABLE courses DROP FOREIGN KEY fk_course_curriculum;
ALTER TABLE courses DROP COLUMN curriculum_id;
DROP TABLE curricula;
