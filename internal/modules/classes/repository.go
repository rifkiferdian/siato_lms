package classes

import "gorm.io/gorm"

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

const classSelect = `SELECT cc.id,c.code course_code,c.name course_name,cc.name class_name,COALESCE(cc.description,'') description,cc.status,c.credits,cc.capacity,cc.max_meetings,
CONCAT(s.name,' ',ay.name) semester_name,ay.name academic_year,
COALESCE((SELECT GROUP_CONCAT(DISTINCT u.name ORDER BY cl.is_coordinator DESC,u.name SEPARATOR ', ') FROM class_lecturers cl JOIN lecturers l ON l.id=cl.lecturer_id JOIN users u ON u.id=l.user_id WHERE cl.class_id=cc.id),'Belum ada dosen') lecturer_names,
(SELECT COUNT(*) FROM class_students cs WHERE cs.class_id=cc.id) student_count,
(SELECT COUNT(*) FROM class_lecturers cl WHERE cl.class_id=cc.id) lecturer_count,
(SELECT COUNT(*) FROM class_schedules sch WHERE sch.class_id=cc.id) schedule_count,
COALESCE((SELECT GROUP_CONCAT(CONCAT(CASE sch.day_of_week WHEN 1 THEN 'Sen' WHEN 2 THEN 'Sel' WHEN 3 THEN 'Rab' WHEN 4 THEN 'Kam' WHEN 5 THEN 'Jum' WHEN 6 THEN 'Sab' ELSE 'Min' END,', ',TIME_FORMAT(sch.start_time,'%H:%i')) ORDER BY sch.day_of_week,sch.start_time SEPARATOR ' · ') FROM class_schedules sch WHERE sch.class_id=cc.id),'Belum dijadwalkan') schedule_summary
FROM course_classes cc JOIN courses c ON c.id=cc.course_id JOIN semesters s ON s.id=cc.semester_id JOIN academic_years ay ON ay.id=s.academic_year_id`

func (r *Repository) List(portal string, userID uint64, search, status string) ([]ClassRow, error) {
	q := classSelect + " WHERE 1=1"
	args := []any{}
	if portal == "lecturer" {
		q += " AND EXISTS(SELECT 1 FROM class_lecturers x JOIN lecturers lx ON lx.id=x.lecturer_id WHERE x.class_id=cc.id AND lx.user_id=?)"
		args = append(args, userID)
	} else if portal == "student" {
		q += " AND EXISTS(SELECT 1 FROM class_students x JOIN students sx ON sx.id=x.student_id WHERE x.class_id=cc.id AND sx.user_id=?)"
		args = append(args, userID)
	}
	if search != "" {
		q += " AND (c.name LIKE ? OR c.code LIKE ? OR cc.name LIKE ?)"
		like := "%" + search + "%"
		args = append(args, like, like, like)
	}
	if status != "" {
		q += " AND cc.status=?"
		args = append(args, status)
	}
	q += " ORDER BY FIELD(cc.status,'active','draft','completed','cancelled'),c.name,cc.name"
	var rows []ClassRow
	e := r.db.Raw(q, args...).Scan(&rows).Error
	return rows, e
}
func (r *Repository) Detail(id, userID uint64, portal string) (ClassDetail, error) {
	q := classSelect + " WHERE cc.id=?"
	args := []any{id}
	if portal == "lecturer" {
		q += " AND EXISTS(SELECT 1 FROM class_lecturers x JOIN lecturers lx ON lx.id=x.lecturer_id WHERE x.class_id=cc.id AND lx.user_id=?)"
		args = append(args, userID)
	} else if portal == "student" {
		q += " AND EXISTS(SELECT 1 FROM class_students x JOIN students sx ON sx.id=x.student_id WHERE x.class_id=cc.id AND sx.user_id=?)"
		args = append(args, userID)
	}
	var d ClassDetail
	e := r.db.Raw(q, args...).Take(&d).Error
	return d, e
}
func (r *Repository) Students(classID uint64) ([]StudentRow, []StudentOption, error) {
	var students []StudentRow
	err := r.db.Raw(`SELECT s.id,u.name,u.email,s.student_number,sp.name program_name
FROM class_students cs JOIN students s ON s.id=cs.student_id JOIN users u ON u.id=s.user_id
JOIN study_programs sp ON sp.id=s.study_program_id WHERE cs.class_id=? ORDER BY u.name`, classID).Scan(&students).Error
	if err != nil {
		return nil, nil, err
	}
	var available []StudentOption
	err = r.db.Raw(`SELECT s.id,s.entry_year,u.name,s.student_number,sp.name program_name
FROM students s JOIN users u ON u.id=s.user_id JOIN study_programs sp ON sp.id=s.study_program_id
WHERE u.is_active=1 AND NOT EXISTS(SELECT 1 FROM class_students cs WHERE cs.class_id=? AND cs.student_id=s.id)
ORDER BY u.name`, classID).Scan(&available).Error
	return students, available, err
}
func (r *Repository) AddStudents(classID uint64, studentIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var capacity, enrolled int
		if err := tx.Raw("SELECT capacity FROM course_classes WHERE id=? FOR UPDATE", classID).Scan(&capacity).Error; err != nil {
			return err
		}
		if capacity == 0 {
			return gormErrNotFound
		}
		if err := tx.Raw("SELECT COUNT(*) FROM class_students WHERE class_id=?", classID).Scan(&enrolled).Error; err != nil {
			return err
		}
		if enrolled+len(studentIDs) > capacity {
			return ErrClassFull
		}
		var valid int64
		if err := tx.Table("students s").Joins("JOIN users u ON u.id=s.user_id").Where("s.id IN ? AND u.is_active=1", studentIDs).Count(&valid).Error; err != nil {
			return err
		}
		if int(valid) != len(studentIDs) {
			return ErrInvalidStudent
		}
		for _, studentID := range studentIDs {
			result := tx.Exec("INSERT IGNORE INTO class_students(class_id,student_id) VALUES(?,?)", classID, studentID)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrStudentAlreadyEnrolled
			}
		}
		return nil
	})
}
func (r *Repository) RemoveStudent(classID, studentID uint64) error {
	result := r.db.Exec("DELETE FROM class_students WHERE class_id=? AND student_id=?", classID, studentID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrEnrollmentNotFound
	}
	return nil
}
func (r *Repository) Options() ([]Option, []Option, error) {
	var courses, semesters []Option
	e := r.db.Raw("SELECT id,CONCAT(code,' · ',name) label FROM courses ORDER BY name").Scan(&courses).Error
	if e != nil {
		return nil, nil, e
	}
	e = r.db.Raw("SELECT s.id,CONCAT(s.name,' · ',ay.name) label FROM semesters s JOIN academic_years ay ON ay.id=s.academic_year_id ORDER BY s.starts_on DESC").Scan(&semesters).Error
	return courses, semesters, e
}
func (r *Repository) Create(in CreateInput) error {
	return r.db.Exec("INSERT INTO course_classes(course_id,semester_id,name,description,capacity,max_meetings,status) VALUES(?,?,?,?,?,?,?)", in.CourseID, in.SemesterID, in.Name, in.Description, in.Capacity, in.MaxMeetings, in.Status).Error
}
func (r *Repository) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if e := tx.Table("class_students").Where("class_id=?", id).Count(&count).Error; e != nil {
			return e
		}
		if count > 0 {
			return ErrClassHasStudents
		}
		return tx.Delete(&courseClass{}, id).Error
	})
}

type courseClass struct{ ID uint64 }

func (courseClass) TableName() string { return "course_classes" }
