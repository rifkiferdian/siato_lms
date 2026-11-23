ALTER TABLE class_lecturer_components ADD COLUMN credits DECIMAL(4,2) NOT NULL DEFAULT 0 AFTER component_id;
UPDATE class_lecturer_components clc
JOIN course_classes cc ON cc.id=clc.class_id
JOIN course_credit_components ccc ON ccc.course_id=cc.course_id AND ccc.component_id=clc.component_id
SET clc.credits=ccc.credits/(SELECT COUNT(*) FROM class_lecturer_components x WHERE x.class_id=clc.class_id AND x.component_id=clc.component_id);
