CREATE TABLE curricula (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
 study_program_id BIGINT UNSIGNED NOT NULL,
 code VARCHAR(50) NOT NULL,
 name VARCHAR(150) NOT NULL,
 effective_year SMALLINT UNSIGNED NOT NULL,
 is_active TINYINT(1) NOT NULL DEFAULT 1,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_curriculum_program_code(study_program_id,code),
 CONSTRAINT fk_curriculum_program FOREIGN KEY(study_program_id) REFERENCES study_programs(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
INSERT INTO curricula(study_program_id,code,name,effective_year,is_active)
SELECT id,'DEFAULT','Kurikulum Awal',YEAR(CURDATE()),1 FROM study_programs;
ALTER TABLE courses ADD COLUMN curriculum_id BIGINT UNSIGNED NULL AFTER study_program_id;
UPDATE courses c JOIN curricula k ON k.study_program_id=c.study_program_id AND k.code='DEFAULT' SET c.curriculum_id=k.id;
ALTER TABLE courses MODIFY curriculum_id BIGINT UNSIGNED NOT NULL;
ALTER TABLE courses ADD CONSTRAINT fk_course_curriculum FOREIGN KEY(curriculum_id) REFERENCES curricula(id);
ALTER TABLE courses ADD COLUMN recommended_semester TINYINT UNSIGNED NOT NULL DEFAULT 1 AFTER credits;
ALTER TABLE courses ADD COLUMN course_type ENUM('required','elective') NOT NULL DEFAULT 'required' AFTER recommended_semester;
ALTER TABLE courses ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 1 AFTER course_type;
DROP INDEX uq_course_code ON courses;
CREATE UNIQUE INDEX uq_course_curriculum_code ON courses(curriculum_id,code);
