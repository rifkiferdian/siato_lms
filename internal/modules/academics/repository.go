package academics

import (
	"errors"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List() ([]AcademicYear, []Semester, []StudyProgram, []Course, error) {
	var years []AcademicYear
	if err := r.db.Raw(`SELECT y.id,y.name,DATE_FORMAT(y.starts_on,'%Y-%m-%d') starts_on,DATE_FORMAT(y.ends_on,'%Y-%m-%d') ends_on,y.is_active,(SELECT COUNT(*) FROM semesters s WHERE s.academic_year_id=y.id) semester_count FROM academic_years y ORDER BY y.starts_on DESC`).Scan(&years).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	var semesters []Semester
	if err := r.db.Raw(`SELECT s.id,s.academic_year_id,s.name,s.code,y.name academic_year_name,DATE_FORMAT(s.starts_on,'%Y-%m-%d') starts_on,DATE_FORMAT(s.ends_on,'%Y-%m-%d') ends_on,s.is_active,(SELECT COUNT(*) FROM course_classes c WHERE c.semester_id=s.id) class_count FROM semesters s JOIN academic_years y ON y.id=s.academic_year_id ORDER BY s.starts_on DESC`).Scan(&semesters).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	var programs []StudyProgram
	if err := r.db.Raw(`SELECT p.id,p.code,p.name,p.degree,(SELECT COUNT(*) FROM students s WHERE s.study_program_id=p.id) student_count,(SELECT COUNT(*) FROM lecturers l WHERE l.study_program_id=p.id) lecturer_count,(SELECT COUNT(*) FROM courses c WHERE c.study_program_id=p.id) course_count FROM study_programs p ORDER BY p.name`).Scan(&programs).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	var courses []Course
	if err := r.db.Raw(`SELECT c.id,c.study_program_id,c.code,c.name,c.credits,p.name program_name,(SELECT COUNT(*) FROM course_classes cc WHERE cc.course_id=c.id) class_count FROM courses c JOIN study_programs p ON p.id=c.study_program_id ORDER BY c.code`).Scan(&courses).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	return years, semesters, programs, courses, nil
}

func (r *Repository) CreateYear(in AcademicYearInput) error {
	return r.db.Exec("INSERT INTO academic_years(name,starts_on,ends_on,is_active) VALUES(?,?,?,0)", in.Name, in.StartsOn, in.EndsOn).Error
}
func (r *Repository) UpdateYear(id uint64, in AcademicYearInput) error {
	return affected(r.db.Exec("UPDATE academic_years SET name=?,starts_on=?,ends_on=? WHERE id=?", in.Name, in.StartsOn, in.EndsOn, id), "tahun akademik")
}
func (r *Repository) DeleteYear(id uint64) error {
	return affected(r.db.Exec("DELETE FROM academic_years WHERE id=?", id), "tahun akademik")
}
func (r *Repository) CreateSemester(in SemesterInput) error {
	return r.db.Exec("INSERT INTO semesters(academic_year_id,name,code,starts_on,ends_on,is_active) VALUES(?,?,?,?,?,0)", in.AcademicYearID, in.Name, in.Code, in.StartsOn, in.EndsOn).Error
}
func (r *Repository) UpdateSemester(id uint64, in SemesterInput) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var active bool
		if err := tx.Raw("SELECT is_active FROM semesters WHERE id=?", id).Scan(&active).Error; err != nil {
			return err
		}
		if err := affected(tx.Exec("UPDATE semesters SET academic_year_id=?,name=?,code=?,starts_on=?,ends_on=? WHERE id=?", in.AcademicYearID, in.Name, in.Code, in.StartsOn, in.EndsOn, id), "semester"); err != nil {
			return err
		}
		if active {
			if err := tx.Exec("UPDATE academic_years SET is_active=0").Error; err != nil {
				return err
			}
			return tx.Exec("UPDATE academic_years SET is_active=1 WHERE id=?", in.AcademicYearID).Error
		}
		return nil
	})
}
func (r *Repository) DeleteSemester(id uint64) error {
	var active bool
	if err := r.db.Raw("SELECT is_active FROM semesters WHERE id=?", id).Scan(&active).Error; err != nil {
		return err
	}
	if active {
		return errors.New("semester aktif tidak dapat dihapus; aktifkan semester lain terlebih dahulu")
	}
	return affected(r.db.Exec("DELETE FROM semesters WHERE id=?", id), "semester")
}
func (r *Repository) ActivateSemester(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var yearID uint64
		if err := tx.Raw("SELECT academic_year_id FROM semesters WHERE id=?", id).Scan(&yearID).Error; err != nil {
			return err
		}
		if yearID == 0 {
			return errors.New("semester tidak ditemukan")
		}
		if err := tx.Exec("UPDATE semesters SET is_active=0").Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE academic_years SET is_active=0").Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE semesters SET is_active=1 WHERE id=?", id).Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE academic_years SET is_active=1 WHERE id=?", yearID).Error
	})
}
func (r *Repository) CreateProgram(in StudyProgramInput) error {
	return r.db.Exec("INSERT INTO study_programs(code,name,degree) VALUES(?,?,?)", in.Code, in.Name, in.Degree).Error
}
func (r *Repository) UpdateProgram(id uint64, in StudyProgramInput) error {
	return affected(r.db.Exec("UPDATE study_programs SET code=?,name=?,degree=? WHERE id=?", in.Code, in.Name, in.Degree, id), "program studi")
}
func (r *Repository) DeleteProgram(id uint64) error {
	return affected(r.db.Exec("DELETE FROM study_programs WHERE id=?", id), "program studi")
}
func (r *Repository) CreateCourse(in CourseInput) error {
	return r.db.Exec("INSERT INTO courses(study_program_id,code,name,credits) VALUES(?,?,?,?)", in.StudyProgramID, in.Code, in.Name, in.Credits).Error
}
func (r *Repository) UpdateCourse(id uint64, in CourseInput) error {
	return affected(r.db.Exec("UPDATE courses SET study_program_id=?,code=?,name=?,credits=? WHERE id=?", in.StudyProgramID, in.Code, in.Name, in.Credits, id), "mata kuliah")
}
func (r *Repository) DeleteCourse(id uint64) error {
	return affected(r.db.Exec("DELETE FROM courses WHERE id=?", id), "mata kuliah")
}
func affected(result *gorm.DB, label string) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New(label + " tidak ditemukan")
	}
	return nil
}
