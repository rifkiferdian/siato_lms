package coursecatalog

import (
	"errors"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db} }
func (r *Repository) Page(q string, programID, curriculumID uint64, page int) (PageData, error) {
	const perPage = 20
	where := " WHERE 1=1"
	args := []any{}
	if q != "" {
		where += " AND (c.code LIKE ? OR c.name LIKE ?)"
		like := "%" + q + "%"
		args = append(args, like, like)
	}
	if programID > 0 {
		where += " AND k.study_program_id=?"
		args = append(args, programID)
	}
	if curriculumID > 0 {
		where += " AND c.curriculum_id=?"
		args = append(args, curriculumID)
	}
	base := " FROM courses c JOIN curricula k ON k.id=c.curriculum_id JOIN study_programs p ON p.id=k.study_program_id"
	var total int
	if e := r.db.Raw("SELECT COUNT(*)"+base+where, args...).Scan(&total).Error; e != nil {
		return PageData{}, e
	}
	pages := (total + perPage - 1) / perPage
	if pages < 1 {
		pages = 1
	}
	if page < 1 {
		page = 1
	}
	if page > pages {
		page = pages
	}
	var rows []Course
	qa := append(append([]any{}, args...), perPage, (page-1)*perPage)
	e := r.db.Raw(`SELECT c.id,c.curriculum_id,c.code,c.name,c.credits,c.recommended_semester,c.course_type,c.is_active,p.name program_name,k.name curriculum_name,
		COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='T'),0) theory_credits,
		COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='P'),0) practicum_credits,
		COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='PR'),0) practice_credits,
		COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='L'),0) field_credits`+base+where+" ORDER BY c.code LIMIT ? OFFSET ?", qa...).Scan(&rows).Error
	for i := range rows {
		rows[i].RowNumber = (page-1)*perPage + i + 1
	}
	var curricula []Curriculum
	_ = r.db.Raw("SELECT k.id,k.study_program_id,k.code,k.name,k.effective_year,k.is_active,p.name program_name FROM curricula k JOIN study_programs p ON p.id=k.study_program_id ORDER BY p.name,k.effective_year DESC").Scan(&curricula).Error
	var programs []Option
	_ = r.db.Raw("SELECT id,CONCAT(code,' - ',name) label FROM study_programs ORDER BY name").Scan(&programs).Error
	return PageData{Search: q, ProgramID: programID, CurriculumID: curriculumID, Page: page, PrevPage: page - 1, NextPage: page + 1, TotalPages: pages, Total: total, Courses: rows, Curricula: curricula, Programs: programs}, e
}
func (r *Repository) SaveCourse(id uint64, in CourseInput) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if id == 0 {
			res := tx.Exec("INSERT INTO courses(study_program_id,curriculum_id,code,name,credits,recommended_semester,course_type,is_active) SELECT study_program_id,id,?,?,?,?,?,? FROM curricula WHERE id=?", in.Code, in.Name, in.Credits, in.RecommendedSemester, in.CourseType, in.IsActive, in.CurriculumID)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return errors.New("kurikulum tidak ditemukan")
			}
			if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error; err != nil {
				return err
			}
		} else if err := tx.Exec("UPDATE courses c JOIN curricula k ON k.id=? SET c.study_program_id=k.study_program_id,c.curriculum_id=k.id,c.code=?,c.name=?,c.credits=?,c.recommended_semester=?,c.course_type=?,c.is_active=? WHERE c.id=?", in.CurriculumID, in.Code, in.Name, in.Credits, in.RecommendedSemester, in.CourseType, in.IsActive, id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM course_credit_components WHERE course_id=?", id).Error; err != nil {
			return err
		}
		components := []struct {
			code    string
			credits int
		}{{"T", in.TheoryCredits}, {"P", in.PracticumCredits}, {"PR", in.PracticeCredits}, {"L", in.FieldCredits}}
		for _, component := range components {
			if component.credits == 0 {
				continue
			}
			if err := tx.Exec("INSERT INTO course_credit_components(course_id,component_id,credits) SELECT ?,id,? FROM learning_components WHERE code=?", id, component.credits, component.code).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *Repository) DeleteCourse(id uint64) error {
	return r.db.Exec("DELETE FROM courses WHERE id=?", id).Error
}
func (r *Repository) SaveCurriculum(id uint64, in CurriculumInput) error {
	if id == 0 {
		return r.db.Exec("INSERT INTO curricula(study_program_id,code,name,effective_year,is_active) VALUES(?,?,?,?,1)", in.StudyProgramID, in.Code, in.Name, in.EffectiveYear).Error
	}
	return r.db.Exec("UPDATE curricula SET study_program_id=?,code=?,name=?,effective_year=? WHERE id=?", in.StudyProgramID, in.Code, in.Name, in.EffectiveYear, id).Error
}
