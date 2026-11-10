DROP INDEX idx_class_status_semester ON course_classes;
ALTER TABLE course_classes DROP COLUMN status;
ALTER TABLE course_classes DROP COLUMN description;
