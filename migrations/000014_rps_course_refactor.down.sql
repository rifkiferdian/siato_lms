ALTER TABLE course_rps ADD COLUMN class_id BIGINT UNSIGNED NULL AFTER course_id;
UPDATE course_rps r JOIN course_classes cc ON cc.rps_id=r.id SET r.class_id=cc.id;
ALTER TABLE course_classes DROP FOREIGN KEY fk_class_rps;
ALTER TABLE course_classes DROP COLUMN rps_id;
ALTER TABLE course_rps DROP FOREIGN KEY fk_rps_course;
ALTER TABLE course_rps DROP INDEX uq_rps_course_version;
ALTER TABLE course_rps DROP COLUMN course_id;
ALTER TABLE course_rps ADD CONSTRAINT fk_rps_class FOREIGN KEY(class_id) REFERENCES course_classes(id) ON DELETE CASCADE;
