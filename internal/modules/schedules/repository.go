package schedules

import "gorm.io/gorm"

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

const selectSchedule = `SELECT s.id,s.class_id,c.code course_code,c.name course_name,cc.name class_name,c.credits,
s.day_of_week, CASE s.day_of_week WHEN 1 THEN 'Senin' WHEN 2 THEN 'Selasa' WHEN 3 THEN 'Rabu' WHEN 4 THEN 'Kamis' WHEN 5 THEN 'Jumat' WHEN 6 THEN 'Sabtu' ELSE 'Minggu' END day_name,
TIME_FORMAT(s.start_time,'%H:%i') start_time,TIME_FORMAT(s.end_time,'%H:%i') end_time,s.room,s.delivery_mode,COALESCE(s.meeting_url,'') meeting_url,
COALESCE(GROUP_CONCAT(DISTINCT u.name ORDER BY cl.is_coordinator DESC,u.name SEPARATOR ', '),'Belum ditentukan') lecturer_names
FROM class_schedules s JOIN course_classes cc ON cc.id=s.class_id JOIN courses c ON c.id=cc.course_id
LEFT JOIN class_lecturers cl ON cl.class_id=cc.id LEFT JOIN lecturers l ON l.id=cl.lecturer_id LEFT JOIN users u ON u.id=l.user_id`

func (r *Repository) List(portal string, userID uint64, day int) ([]Row, error) {
	query := selectSchedule
	args := []any{}
	switch portal {
	case "lecturer":
		query += " WHERE EXISTS (SELECT 1 FROM class_lecturers x JOIN lecturers lx ON lx.id=x.lecturer_id WHERE x.class_id=cc.id AND lx.user_id=?)"
		args = append(args, userID)
	case "student":
		query += " WHERE EXISTS (SELECT 1 FROM class_students x JOIN students sx ON sx.id=x.student_id WHERE x.class_id=cc.id AND sx.user_id=?)"
		args = append(args, userID)
	default:
		query += " WHERE 1=1"
	}
	if day > 0 {
		query += " AND s.day_of_week=?"
		args = append(args, day)
	}
	query += " GROUP BY s.id,c.code,c.name,cc.name,c.credits,s.day_of_week,s.start_time,s.end_time,s.room,s.delivery_mode,s.meeting_url ORDER BY s.day_of_week,s.start_time,c.name"
	var rows []Row
	err := r.db.Raw(query, args...).Scan(&rows).Error
	return rows, err
}
func (r *Repository) Classes() ([]ClassOption, error) {
	var out []ClassOption
	err := r.db.Raw(`SELECT cc.id,CONCAT(c.code,' · ',c.name,' · Kelas ',cc.name) label FROM course_classes cc JOIN courses c ON c.id=cc.course_id JOIN semesters s ON s.id=cc.semester_id WHERE s.is_active=1 ORDER BY c.name,cc.name`).Scan(&out).Error
	return out, err
}
func (r *Repository) HasConflict(in CreateInput) (bool, error) {
	var n int64
	err := r.db.Raw(`SELECT COUNT(*) FROM class_schedules WHERE day_of_week=? AND (class_id=? OR LOWER(room)=LOWER(?)) AND start_time < ? AND end_time > ?`, in.DayOfWeek, in.ClassID, in.Room, in.EndTime, in.StartTime).Scan(&n).Error
	return n > 0, err
}
func (r *Repository) Create(in CreateInput) error {
	return r.db.Exec(`INSERT INTO class_schedules(class_id,day_of_week,start_time,end_time,room,delivery_mode,meeting_url,effective_from,effective_until) VALUES(?,?,?,?,?,?,NULLIF(?,''),NULLIF(?,''),NULLIF(?,''))`, in.ClassID, in.DayOfWeek, in.StartTime, in.EndTime, in.Room, in.DeliveryMode, in.MeetingURL, in.EffectiveFrom, in.EffectiveUntil).Error
}
func (r *Repository) Delete(id uint64) error { return r.db.Delete(&Schedule{}, id).Error }
