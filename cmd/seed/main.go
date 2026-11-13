package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"os"
	"siato_lms/internal/config"
	"siato_lms/internal/database"
)

func main() {
	c, e := config.Load()
	if e != nil {
		log.Fatal(e)
	}
	db, e := database.Open(c, slog.Default())
	if e != nil {
		log.Fatal(e)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	tx := db.Begin()
	roles := []string{"super_admin", "academic_admin", "lecturer", "student"}
	for _, r := range roles {
		tx.Exec("INSERT IGNORE INTO roles(name,guard_name) VALUES(?, 'web')", r)
	}
	perms := []string{"admin.dashboard", "lecturer.dashboard", "student.dashboard", "classes.view", "classes.manage", "materials.view", "materials.create", "materials.update", "assignments.create", "assignments.submit", "assignments.grade", "attendance.open", "attendance.scan", "attendance.manage", "messaging.send", "announcements.create", "reports.view", "users.view", "users.manage", "roles.view", "roles.manage"}
	for _, p := range perms {
		tx.Exec("INSERT IGNORE INTO permissions(name,guard_name) VALUES(?, 'web')", p)
	}
	users := [][]string{{"Super Admin", "admin@siato.test", "super_admin"}, {"Admin Akademik", "akademik@siato.test", "academic_admin"}, {"Dosen Demo", "dosen@siato.test", "lecturer"}, {"Mahasiswa Demo", "mahasiswa@siato.test", "student"}}
	for _, u := range users {
		tx.Exec("INSERT INTO users(name,email,password_hash,is_active) VALUES(?,?,?,1) ON DUPLICATE KEY UPDATE name=VALUES(name)", u[0], u[1], string(hash))
		tx.Exec("INSERT IGNORE INTO model_has_roles(role_id,user_id) SELECT r.id,u.id FROM roles r JOIN users u ON u.email=? WHERE r.name=?", u[1], u[2])
	}
	// Academic fixtures for the schedule module. INSERT IGNORE keeps the seeder repeatable.
	tx.Exec("INSERT IGNORE INTO academic_years(name,starts_on,ends_on,is_active,external_id) VALUES('2026/2027','2026-08-01','2027-07-31',1,'AY-2026')")
	tx.Exec("INSERT IGNORE INTO semesters(academic_year_id,name,code,starts_on,ends_on,is_active) SELECT id,'Ganjil','2026-GANJIL','2026-08-18','2027-01-15',1 FROM academic_years WHERE name='2026/2027'")
	tx.Exec("INSERT IGNORE INTO study_programs(code,name,degree,external_id) VALUES('IF','Informatika','S1','PRODI-IF')")
	tx.Exec("INSERT IGNORE INTO lecturers(user_id,study_program_id,employee_number,external_id) SELECT u.id,p.id,'DSN001','LECT-001' FROM users u JOIN study_programs p ON p.code='IF' WHERE u.email='dosen@siato.test'")
	tx.Exec("INSERT IGNORE INTO students(user_id,study_program_id,student_number,entry_year,external_id) SELECT u.id,p.id,'20260001',2026,'STU-001' FROM users u JOIN study_programs p ON p.code='IF' WHERE u.email='mahasiswa@siato.test'")
	dummyStudents := []struct {
		name, email, number string
	}{
		{"Andi Pratama", "andi.pratama@siato.test", "20260002"},
		{"Bunga Lestari", "bunga.lestari@siato.test", "20260003"},
		{"Cahya Nugraha", "cahya.nugraha@siato.test", "20260004"},
		{"Dinda Maharani", "dinda.maharani@siato.test", "20260005"},
		{"Eko Saputra", "eko.saputra@siato.test", "20260006"},
		{"Farah Azzahra", "farah.azzahra@siato.test", "20260007"},
		{"Galih Ramadhan", "galih.ramadhan@siato.test", "20260008"},
		{"Hana Permata", "hana.permata@siato.test", "20260009"},
		{"Iqbal Maulana", "iqbal.maulana@siato.test", "20260010"},
	}
	for index, student := range dummyStudents {
		tx.Exec("INSERT INTO users(name,email,password_hash,is_active) VALUES(?,?,?,1) ON DUPLICATE KEY UPDATE name=VALUES(name)", student.name, student.email, string(hash))
		tx.Exec("INSERT IGNORE INTO model_has_roles(role_id,user_id) SELECT r.id,u.id FROM roles r JOIN users u ON u.email=? WHERE r.name='student'", student.email)
		tx.Exec("INSERT IGNORE INTO students(user_id,study_program_id,student_number,entry_year,external_id) SELECT u.id,p.id,?,2026,? FROM users u JOIN study_programs p ON p.code='IF' WHERE u.email=?", student.number, fmt.Sprintf("STU-%03d", index+2), student.email)
	}
	tx.Exec("INSERT IGNORE INTO courses(study_program_id,code,name,credits,external_id) SELECT id,'IF101','Algoritma dan Pemrograman',3,'COURSE-IF101' FROM study_programs WHERE code='IF'")
	tx.Exec("INSERT IGNORE INTO courses(study_program_id,code,name,credits,external_id) SELECT id,'IF203','Basis Data',3,'COURSE-IF203' FROM study_programs WHERE code='IF'")
	tx.Exec("INSERT IGNORE INTO courses(study_program_id,code,name,credits,external_id) SELECT id,'IF305','Rekayasa Perangkat Lunak',3,'COURSE-IF305' FROM study_programs WHERE code='IF'")
	tx.Exec("INSERT IGNORE INTO course_classes(course_id,semester_id,name,capacity,max_meetings,external_id) SELECT c.id,s.id,'A',40,16,CONCAT('CLASS-',c.code,'-A') FROM courses c JOIN semesters s ON s.code='2026-GANJIL' WHERE c.code IN ('IF101','IF203','IF305')")
	tx.Exec("INSERT IGNORE INTO class_lecturers(class_id,lecturer_id,is_coordinator) SELECT cc.id,l.id,1 FROM course_classes cc JOIN courses c ON c.id=cc.course_id JOIN lecturers l JOIN users u ON u.id=l.user_id WHERE c.code IN ('IF101','IF203','IF305') AND u.email='dosen@siato.test'")
	tx.Exec("INSERT IGNORE INTO class_students(class_id,student_id) SELECT cc.id,s.id FROM course_classes cc JOIN courses c ON c.id=cc.course_id JOIN students s JOIN users u ON u.id=s.user_id WHERE c.code IN ('IF101','IF203','IF305') AND u.email='mahasiswa@siato.test'")
	tx.Exec("INSERT IGNORE INTO class_schedules(class_id,day_of_week,start_time,end_time,room,delivery_mode) SELECT cc.id,1,'08:00','10:30','Lab Komputer 1','onsite' FROM course_classes cc JOIN courses c ON c.id=cc.course_id WHERE c.code='IF101'")
	tx.Exec("INSERT IGNORE INTO class_schedules(class_id,day_of_week,start_time,end_time,room,delivery_mode) SELECT cc.id,3,'10:30','13:00','Ruang 3.02','hybrid' FROM course_classes cc JOIN courses c ON c.id=cc.course_id WHERE c.code='IF203'")
	tx.Exec("INSERT IGNORE INTO class_schedules(class_id,day_of_week,start_time,end_time,room,delivery_mode) SELECT cc.id,5,'13:30','16:00','Ruang 2.04','onsite' FROM course_classes cc JOIN courses c ON c.id=cc.course_id WHERE c.code='IF305'")
	tx.Exec("INSERT IGNORE INTO role_has_permissions(permission_id,role_id) SELECT p.id,r.id FROM permissions p JOIN roles r WHERE r.name='super_admin'")
	tx.Exec("INSERT IGNORE INTO role_has_permissions(permission_id,role_id) SELECT p.id,r.id FROM permissions p JOIN roles r WHERE (r.name='academic_admin' AND p.name IN ('admin.dashboard','classes.view','classes.manage','reports.view')) OR (r.name='lecturer' AND p.name IN ('lecturer.dashboard','classes.view','materials.view','materials.create','materials.update','assignments.create','assignments.grade','attendance.open','attendance.manage','messaging.send','announcements.create','reports.view')) OR (r.name='student' AND p.name IN ('student.dashboard','classes.view','materials.view','assignments.submit','attendance.scan','messaging.send'))")
	if e = tx.Commit().Error; e != nil {
		log.Fatal(e)
	}
	log.New(os.Stdout, "", 0).Println("seed complete; demo password: Password123! (development only)")
}
