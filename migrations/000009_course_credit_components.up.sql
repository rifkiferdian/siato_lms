CREATE TABLE learning_components (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
 code VARCHAR(10) NOT NULL,
 name VARCHAR(80) NOT NULL,
 sort_order TINYINT UNSIGNED NOT NULL DEFAULT 0,
 is_active TINYINT(1) NOT NULL DEFAULT 1,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_learning_component_code(code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO learning_components(code,name,sort_order,is_active) VALUES
('T','Teori',1,1),('P','Praktikum',2,1),('PR','Praktik',3,1),('L','Lapangan',4,1);

CREATE TABLE course_credit_components (
 course_id BIGINT UNSIGNED NOT NULL,
 component_id BIGINT UNSIGNED NOT NULL,
 credits TINYINT UNSIGNED NOT NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 PRIMARY KEY(course_id,component_id),
 CONSTRAINT fk_course_credit_course FOREIGN KEY(course_id) REFERENCES courses(id) ON DELETE CASCADE,
 CONSTRAINT fk_course_credit_component FOREIGN KEY(component_id) REFERENCES learning_components(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO course_credit_components(course_id,component_id,credits)
SELECT c.id,lc.id,c.credits FROM courses c JOIN learning_components lc ON lc.code='T' WHERE c.credits>0;
