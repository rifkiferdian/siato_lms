CREATE TABLE academic_years (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, name VARCHAR(20) NOT NULL, starts_on DATE NOT NULL, ends_on DATE NOT NULL, is_active TINYINT(1) NOT NULL DEFAULT 0,
 external_id VARCHAR(100) NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_academic_year_name(name), UNIQUE KEY uq_academic_year_external(external_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE semesters (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, academic_year_id BIGINT UNSIGNED NOT NULL, name VARCHAR(30) NOT NULL, code VARCHAR(30) NOT NULL, starts_on DATE NOT NULL, ends_on DATE NOT NULL, is_active TINYINT(1) NOT NULL DEFAULT 0,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_semester_code(code), KEY idx_semester_year(academic_year_id), CONSTRAINT fk_semester_year FOREIGN KEY(academic_year_id) REFERENCES academic_years(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE study_programs (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, code VARCHAR(30) NOT NULL, name VARCHAR(150) NOT NULL, degree VARCHAR(20) NOT NULL, external_id VARCHAR(100) NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_program_code(code), UNIQUE KEY uq_program_external(external_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE lecturers (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, user_id BIGINT UNSIGNED NOT NULL, study_program_id BIGINT UNSIGNED NULL, employee_number VARCHAR(50) NOT NULL, external_id VARCHAR(100) NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_lecturer_user(user_id), UNIQUE KEY uq_lecturer_number(employee_number), UNIQUE KEY uq_lecturer_external(external_id),
 CONSTRAINT fk_lecturer_user FOREIGN KEY(user_id) REFERENCES users(id), CONSTRAINT fk_lecturer_program FOREIGN KEY(study_program_id) REFERENCES study_programs(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE students (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, user_id BIGINT UNSIGNED NOT NULL, study_program_id BIGINT UNSIGNED NOT NULL, student_number VARCHAR(50) NOT NULL, entry_year SMALLINT UNSIGNED NOT NULL, external_id VARCHAR(100) NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_student_user(user_id), UNIQUE KEY uq_student_number(student_number), UNIQUE KEY uq_student_external(external_id),
 CONSTRAINT fk_student_user FOREIGN KEY(user_id) REFERENCES users(id), CONSTRAINT fk_student_program FOREIGN KEY(study_program_id) REFERENCES study_programs(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE courses (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, study_program_id BIGINT UNSIGNED NOT NULL, code VARCHAR(30) NOT NULL, name VARCHAR(180) NOT NULL, credits TINYINT UNSIGNED NOT NULL, external_id VARCHAR(100) NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_course_code(code), UNIQUE KEY uq_course_external(external_id), CONSTRAINT fk_course_program FOREIGN KEY(study_program_id) REFERENCES study_programs(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE course_classes (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, course_id BIGINT UNSIGNED NOT NULL, semester_id BIGINT UNSIGNED NOT NULL, name VARCHAR(50) NOT NULL, capacity SMALLINT UNSIGNED NOT NULL DEFAULT 40, max_meetings TINYINT UNSIGNED NOT NULL DEFAULT 16, external_id VARCHAR(100) NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_class_course_semester_name(course_id,semester_id,name), UNIQUE KEY uq_class_external(external_id), KEY idx_class_semester(semester_id),
 CONSTRAINT fk_class_course FOREIGN KEY(course_id) REFERENCES courses(id), CONSTRAINT fk_class_semester FOREIGN KEY(semester_id) REFERENCES semesters(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE class_lecturers (
 class_id BIGINT UNSIGNED NOT NULL, lecturer_id BIGINT UNSIGNED NOT NULL, is_coordinator TINYINT(1) NOT NULL DEFAULT 0, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY(class_id,lecturer_id), KEY idx_cl_lecturer(lecturer_id), CONSTRAINT fk_cl_class FOREIGN KEY(class_id) REFERENCES course_classes(id) ON DELETE CASCADE, CONSTRAINT fk_cl_lecturer FOREIGN KEY(lecturer_id) REFERENCES lecturers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE class_students (
 class_id BIGINT UNSIGNED NOT NULL, student_id BIGINT UNSIGNED NOT NULL, enrolled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY(class_id,student_id), KEY idx_cs_student(student_id), CONSTRAINT fk_cs_class FOREIGN KEY(class_id) REFERENCES course_classes(id) ON DELETE CASCADE, CONSTRAINT fk_cs_student FOREIGN KEY(student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE class_schedules (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY, class_id BIGINT UNSIGNED NOT NULL, day_of_week TINYINT UNSIGNED NOT NULL COMMENT '1=Monday ... 7=Sunday', start_time TIME NOT NULL, end_time TIME NOT NULL,
 room VARCHAR(100) NOT NULL, delivery_mode ENUM('onsite','online','hybrid') NOT NULL DEFAULT 'onsite', meeting_url VARCHAR(500) NULL, effective_from DATE NULL, effective_until DATE NULL,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uq_class_schedule(class_id,day_of_week,start_time), KEY idx_schedule_day_time(day_of_week,start_time),
 CONSTRAINT fk_schedule_class FOREIGN KEY(class_id) REFERENCES course_classes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
INSERT IGNORE INTO permissions(name,guard_name) VALUES ('schedules.view','web'),('schedules.manage','web');
INSERT IGNORE INTO role_has_permissions(permission_id,role_id) SELECT p.id,r.id FROM permissions p JOIN roles r WHERE r.name='super_admin' AND p.name IN ('schedules.view','schedules.manage');
INSERT IGNORE INTO role_has_permissions(permission_id,role_id) SELECT p.id,r.id FROM permissions p JOIN roles r WHERE r.name='academic_admin' AND p.name IN ('schedules.view','schedules.manage');
INSERT IGNORE INTO role_has_permissions(permission_id,role_id) SELECT p.id,r.id FROM permissions p JOIN roles r WHERE r.name IN ('lecturer','student') AND p.name='schedules.view';
