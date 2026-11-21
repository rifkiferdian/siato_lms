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
func (r *Repository) Students(search, status, source string, page, perPage int) ([]StudentRow, int, int, error) {
	base := ` FROM students s JOIN users u ON u.id=s.user_id JOIN study_programs p ON p.id=s.study_program_id WHERE 1=1`
	where := ""
	args := []any{}
	if search != "" {
		like := "%" + search + "%"
		where += " AND (u.name LIKE ? OR u.email LIKE ? OR s.student_number LIKE ?)"
		args = append(args, like, like, like)
	}
	if status != "" {
		where += " AND s.academic_status=?"
		args = append(args, status)
	}
	if source != "" {
		where += " AND s.data_source=?"
		args = append(args, source)
	}
	var total int
	if e := r.db.Raw("SELECT COUNT(*)"+base+where, args...).Scan(&total).Error; e != nil {
		return nil, 0, page, e
	}
	page = clampPage(page, total, perPage)
	q := `SELECT s.id,s.user_id,u.name,u.email,u.is_active,s.student_number,s.entry_year,s.academic_status,s.data_source,COALESCE(s.external_id,'') external_id,COALESCE(DATE_FORMAT(s.last_synced_at,'%d/%m/%Y %H:%i'),'-') last_synced_at,p.code program_code,p.name program_name,(SELECT COUNT(*) FROM class_students cs WHERE cs.student_id=s.id) class_count` + base + where + " ORDER BY u.name LIMIT ? OFFSET ?"
	queryArgs := append(append([]any{}, args...), perPage, (page-1)*perPage)
	var rows []StudentRow
	e := r.db.Raw(q, queryArgs...).Scan(&rows).Error
	return rows, total, page, e
}
func (r *Repository) Lecturers(search, status, source string, page, perPage int) ([]LecturerRow, int, int, error) {
	base := ` FROM lecturers l JOIN users u ON u.id=l.user_id LEFT JOIN study_programs p ON p.id=l.study_program_id WHERE 1=1`
	where := ""
	args := []any{}
	if search != "" {
		like := "%" + search + "%"
		where += " AND (u.name LIKE ? OR u.email LIKE ? OR l.employee_number LIKE ?)"
		args = append(args, like, like, like)
	}
	if status != "" {
		where += " AND l.academic_status=?"
		args = append(args, status)
	}
	if source != "" {
		where += " AND l.data_source=?"
		args = append(args, source)
	}
	var total int
	if e := r.db.Raw("SELECT COUNT(*)"+base+where, args...).Scan(&total).Error; e != nil {
		return nil, 0, page, e
	}
	page = clampPage(page, total, perPage)
	q := `SELECT l.id,l.user_id,u.name,u.email,u.is_active,l.employee_number,l.academic_status,l.data_source,COALESCE(l.external_id,'') external_id,COALESCE(DATE_FORMAT(l.last_synced_at,'%d/%m/%Y %H:%i'),'-') last_synced_at,COALESCE(p.code,'-') program_code,COALESCE(p.name,'Belum ditentukan') program_name,(SELECT COUNT(*) FROM class_lecturers cl WHERE cl.lecturer_id=l.id) class_count` + base + where + " ORDER BY u.name LIMIT ? OFFSET ?"
	queryArgs := append(append([]any{}, args...), perPage, (page-1)*perPage)
	var rows []LecturerRow
	e := r.db.Raw(q, queryArgs...).Scan(&rows).Error
	return rows, total, page, e
}

func clampPage(page, total, perPage int) int {
	if page < 1 {
		page = 1
	}
	pages := (total + perPage - 1) / perPage
	if pages > 0 && page > pages {
		page = pages
	}
	return page
}
func (r *Repository) Programs() ([]ProgramOption, error) {
	var rows []ProgramOption
	e := r.db.Raw("SELECT id,name label FROM study_programs ORDER BY name").Scan(&rows).Error
	return rows, e
}

func (r *Repository) StudentStats() (PeopleStats, error) {
	var stats PeopleStats
	err := r.db.Raw(`SELECT COUNT(*) total,
COALESCE(SUM(u.is_active=1),0) active,
COALESCE(SUM(u.is_active=0),0) inactive,
COALESCE(SUM(s.data_source='siakad'),0) siakad,
COALESCE(SUM(s.data_source='manual'),0) manual
FROM students s JOIN users u ON u.id=s.user_id`).Scan(&stats).Error
	return stats, err
}

func (r *Repository) LecturerStats() (PeopleStats, error) {
	var stats PeopleStats
	err := r.db.Raw(`SELECT COUNT(*) total,
COALESCE(SUM(u.is_active=1),0) active,
COALESCE(SUM(u.is_active=0),0) inactive,
COALESCE(SUM(l.data_source='siakad'),0) siakad,
COALESCE(SUM(l.data_source='manual'),0) manual
FROM lecturers l JOIN users u ON u.id=l.user_id`).Scan(&stats).Error
	return stats, err
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
