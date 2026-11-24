package rps

import "gorm.io/gorm"

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db} }
func (r *Repository) Course(courseID, userID uint64, portal string) (CourseInfo, bool, error) {
	q := `SELECT c.id,c.code,c.name,c.credits,c.recommended_semester,c.course_type,c.is_active,p.name program_name,k.name curriculum_name,COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='T'),0) theory_credits,COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='P'),0) practicum_credits,COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='PR'),0) practice_credits,COALESCE((SELECT cc.credits FROM course_credit_components cc JOIN learning_components lc ON lc.id=cc.component_id WHERE cc.course_id=c.id AND lc.code='L'),0) field_credits FROM courses c JOIN curricula k ON k.id=c.curriculum_id JOIN study_programs p ON p.id=k.study_program_id WHERE c.id=?`
	var info CourseInfo
	e := r.db.Raw(q, courseID).Take(&info).Error
	if e != nil {
		return info, false, e
	}
	manage := portal == "admin"
	if portal == "lecturer" {
		var n int64
		_ = r.db.Raw(`SELECT COUNT(*) FROM course_classes cc JOIN class_lecturers cl ON cl.class_id=cc.id JOIN lecturers l ON l.id=cl.lecturer_id WHERE cc.course_id=? AND l.user_id=? AND cl.is_coordinator=1`, courseID, userID).Scan(&n).Error
		manage = n > 0
	}
	return info, manage, nil
}
func (r *Repository) Versions(courseID uint64) ([]RPS, error) {
	var rows []RPS
	e := r.db.Raw(`SELECT id,course_id,version,COALESCE(description,'') description,COALESCE(learning_methods,'') learning_methods,COALESCE(main_references,'') main_references,COALESCE(supporting_references,'') supporting_references,status,COALESCE(DATE_FORMAT(activated_at,'%d/%m/%Y %H:%i'),'') activated_at,DATE_FORMAT(updated_at,'%d/%m/%Y %H:%i') updated_at FROM course_rps WHERE course_id=? ORDER BY (status='active') DESC,updated_at DESC`, courseID).Scan(&rows).Error
	return rows, e
}
func (r *Repository) Get(courseID, rpsID uint64) (RPS, []CPMK, error) {
	q := `SELECT id,course_id,version,COALESCE(description,'') description,COALESCE(learning_methods,'') learning_methods,COALESCE(main_references,'') main_references,COALESCE(supporting_references,'') supporting_references,status,COALESCE(DATE_FORMAT(activated_at,'%d/%m/%Y %H:%i'),'') activated_at,DATE_FORMAT(updated_at,'%d/%m/%Y %H:%i') updated_at FROM course_rps WHERE course_id=?`
	args := []any{courseID}
	if rpsID > 0 {
		q += " AND id=?"
		args = append(args, rpsID)
	}
	q += " ORDER BY (status='active') DESC,updated_at DESC LIMIT 1"
	var doc RPS
	e := r.db.Raw(q, args...).Scan(&doc).Error
	if e != nil {
		return doc, nil, e
	}
	var rows []CPMK
	if doc.ID > 0 {
		e = r.db.Raw("SELECT id,code,title,description,sort_order FROM rps_cpmk WHERE rps_id=? ORDER BY sort_order,code", doc.ID).Scan(&rows).Error
	}
	return doc, rows, e
}
func (r *Repository) Save(courseID, rpsID, userID uint64, in RPSInput) error {
	if rpsID == 0 {
		return r.db.Exec(`INSERT INTO course_rps(course_id,version,description,learning_methods,main_references,supporting_references,created_by,updated_by) VALUES(?,?,?,?,?,?,?,?)`, courseID, in.Version, in.Description, in.LearningMethods, in.MainReferences, in.SupportingReferences, userID, userID).Error
	}
	return r.db.Exec(`UPDATE course_rps SET version=?,description=?,learning_methods=?,main_references=?,supporting_references=?,updated_by=?,status='draft',activated_at=NULL WHERE id=? AND course_id=?`, in.Version, in.Description, in.LearningMethods, in.MainReferences, in.SupportingReferences, userID, rpsID, courseID).Error
}
func (r *Repository) Activate(courseID, rpsID, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("UPDATE course_rps SET status='draft',activated_at=NULL WHERE course_id=?", courseID).Error; e != nil {
			return e
		}
		return tx.Exec("UPDATE course_rps SET status='active',activated_at=UTC_TIMESTAMP(),updated_by=? WHERE id=? AND course_id=?", userID, rpsID, courseID).Error
	})
}
func (r *Repository) SaveCPMK(courseID, rpsID, id uint64, in CPMKInput) error {
	if id == 0 {
		return r.db.Exec(`INSERT INTO rps_cpmk(rps_id,code,title,description,sort_order) SELECT id,?,?,?,? FROM course_rps WHERE id=? AND course_id=?`, in.Code, in.Title, in.Description, in.SortOrder, rpsID, courseID).Error
	}
	return r.db.Exec(`UPDATE rps_cpmk c JOIN course_rps r ON r.id=c.rps_id SET c.code=?,c.title=?,c.description=?,c.sort_order=?,r.status='draft',r.activated_at=NULL WHERE c.id=? AND r.id=? AND r.course_id=?`, in.Code, in.Title, in.Description, in.SortOrder, id, rpsID, courseID).Error
}
func (r *Repository) DeleteCPMK(courseID, rpsID, id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec(`DELETE c FROM rps_cpmk c JOIN course_rps r ON r.id=c.rps_id WHERE c.id=? AND r.id=? AND r.course_id=?`, id, rpsID, courseID).Error; e != nil {
			return e
		}
		return tx.Exec("UPDATE course_rps SET status='draft',activated_at=NULL WHERE id=? AND course_id=?", rpsID, courseID).Error
	})
}
