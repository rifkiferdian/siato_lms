package people

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) Students(search, status, source string) ([]StudentRow, error) {
	q := `SELECT s.id,s.user_id,u.name,u.email,u.is_active,s.student_number,s.entry_year,s.academic_status,s.data_source,COALESCE(s.external_id,'') external_id,COALESCE(DATE_FORMAT(s.last_synced_at,'%d/%m/%Y %H:%i'),'-') last_synced_at,p.code program_code,p.name program_name,(SELECT COUNT(*) FROM class_students cs WHERE cs.student_id=s.id) class_count FROM students s JOIN users u ON u.id=s.user_id JOIN study_programs p ON p.id=s.study_program_id WHERE 1=1`
	args := []any{}
	if search != "" {
		like := "%" + search + "%"
		q += " AND (u.name LIKE ? OR u.email LIKE ? OR s.student_number LIKE ?)"
		args = append(args, like, like, like)
	}
	if status != "" {
		q += " AND s.academic_status=?"
		args = append(args, status)
	}
	if source != "" {
		q += " AND s.data_source=?"
		args = append(args, source)
	}
	q += " ORDER BY u.name"
	var rows []StudentRow
	e := r.db.Raw(q, args...).Scan(&rows).Error
	return rows, e
}
func (r *Repository) Lecturers(search, status, source string) ([]LecturerRow, error) {
	q := `SELECT l.id,l.user_id,u.name,u.email,u.is_active,l.employee_number,l.academic_status,l.data_source,COALESCE(l.external_id,'') external_id,COALESCE(DATE_FORMAT(l.last_synced_at,'%d/%m/%Y %H:%i'),'-') last_synced_at,COALESCE(p.code,'-') program_code,COALESCE(p.name,'Belum ditentukan') program_name,(SELECT COUNT(*) FROM class_lecturers cl WHERE cl.lecturer_id=l.id) class_count FROM lecturers l JOIN users u ON u.id=l.user_id LEFT JOIN study_programs p ON p.id=l.study_program_id WHERE 1=1`
	args := []any{}
	if search != "" {
		like := "%" + search + "%"
		q += " AND (u.name LIKE ? OR u.email LIKE ? OR l.employee_number LIKE ?)"
		args = append(args, like, like, like)
	}
	if status != "" {
		q += " AND l.academic_status=?"
		args = append(args, status)
	}
	if source != "" {
		q += " AND l.data_source=?"
		args = append(args, source)
	}
	q += " ORDER BY u.name"
	var rows []LecturerRow
	e := r.db.Raw(q, args...).Scan(&rows).Error
	return rows, e
}
func (r *Repository) Programs() ([]ProgramOption, error) {
	var rows []ProgramOption
	e := r.db.Raw("SELECT id,CONCAT(code,' · ',name) label FROM study_programs ORDER BY name").Scan(&rows).Error
	return rows, e
}
func (r *Repository) CreateStudent(in StudentInput, password string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			return e
		}
		res := tx.Exec("INSERT INTO users(name,email,password_hash,is_active) VALUES(?,?,?,1)", in.Name, strings.ToLower(in.Email), string(hash))
		if res.Error != nil {
			return res.Error
		}
		userID := uint64(0)
		if e = tx.Raw("SELECT LAST_INSERT_ID()").Scan(&userID).Error; e != nil {
			return e
		}
		if e = tx.Exec("INSERT INTO students(user_id,study_program_id,student_number,entry_year,academic_status,external_id,data_source) VALUES(?,?,?,?, 'active',NULLIF(?,''),'manual')", userID, in.StudyProgramID, in.StudentNumber, in.EntryYear, in.ExternalID).Error; e != nil {
			return e
		}
		return tx.Exec("INSERT INTO model_has_roles(role_id,user_id) SELECT id,? FROM roles WHERE name='student'", userID).Error
	})
}
func (r *Repository) CreateLecturer(in LecturerInput, password string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			return e
		}
		res := tx.Exec("INSERT INTO users(name,email,password_hash,is_active) VALUES(?,?,?,1)", in.Name, strings.ToLower(in.Email), string(hash))
		if res.Error != nil {
			return res.Error
		}
		var userID uint64
		if e = tx.Raw("SELECT LAST_INSERT_ID()").Scan(&userID).Error; e != nil {
			return e
		}
		if e = tx.Exec("INSERT INTO lecturers(user_id,study_program_id,employee_number,academic_status,external_id,data_source) VALUES(?,?,?,'active',NULLIF(?,''),'manual')", userID, in.StudyProgramID, in.EmployeeNumber, in.ExternalID).Error; e != nil {
			return e
		}
		return tx.Exec("INSERT INTO model_has_roles(role_id,user_id) SELECT id,? FROM roles WHERE name='lecturer'", userID).Error
	})
}
func (r *Repository) Toggle(userID uint64) error {
	res := r.db.Exec("UPDATE users SET is_active=NOT is_active,updated_at=? WHERE id=?", time.Now().UTC(), userID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("pengguna tidak ditemukan")
	}
	return nil
}
