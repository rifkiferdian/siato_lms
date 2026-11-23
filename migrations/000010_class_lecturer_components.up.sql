CREATE TABLE class_lecturer_components (
 class_id BIGINT UNSIGNED NOT NULL,
 lecturer_id BIGINT UNSIGNED NOT NULL,
 component_id BIGINT UNSIGNED NOT NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY(class_id,lecturer_id,component_id),
 CONSTRAINT fk_clc_class_lecturer FOREIGN KEY(class_id,lecturer_id) REFERENCES class_lecturers(class_id,lecturer_id) ON DELETE CASCADE,
 CONSTRAINT fk_clc_component FOREIGN KEY(component_id) REFERENCES learning_components(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
INSERT INTO class_lecturer_components(class_id,lecturer_id,component_id)
SELECT cl.class_id,cl.lecturer_id,lc.id FROM class_lecturers cl JOIN learning_components lc ON lc.code='T';
