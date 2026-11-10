ALTER TABLE course_classes ADD COLUMN description TEXT NULL AFTER name;
ALTER TABLE course_classes ADD COLUMN status ENUM('draft','active','completed','cancelled') NOT NULL DEFAULT 'active' AFTER max_meetings;
CREATE INDEX idx_class_status_semester ON course_classes(status,semester_id);
