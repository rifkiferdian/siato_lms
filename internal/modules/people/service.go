package people

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{repo: r} }
func (s *Service) Students(q, status, source string) ([]StudentRow, []ProgramOption, error) {
	if !studentStatus(status) {
		status = ""
	}
	if source != "manual" && source != "siakad" {
		source = ""
	}
	rows, e := s.repo.Students(strings.TrimSpace(q), status, source)
	if e != nil {
		return nil, nil, e
	}
	programs, e := s.repo.Programs()
	return rows, programs, e
}
func (s *Service) Lecturers(q, status, source string) ([]LecturerRow, []ProgramOption, error) {
	if !lecturerStatus(status) {
		status = ""
	}
	if source != "manual" && source != "siakad" {
		source = ""
	}
	rows, e := s.repo.Lecturers(strings.TrimSpace(q), status, source)
	if e != nil {
		return nil, nil, e
	}
	programs, e := s.repo.Programs()
	return rows, programs, e
}
func (s *Service) CreateStudent(in StudentInput) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(in.Email)
	in.StudentNumber = strings.TrimSpace(in.StudentNumber)
	if in.Name == "" || in.StudentNumber == "" || in.StudyProgramID == 0 {
		return errors.New("nama, NIM, dan program studi wajib diisi")
	}
	if _, e := mail.ParseAddress(in.Email); e != nil {
		return errors.New("alamat email tidak valid")
	}
	year := time.Now().Year()
	if in.EntryYear < 2000 || in.EntryYear > year+1 {
		return errors.New("tahun masuk tidak valid")
	}
	if len(in.Password) < 8 {
		return errors.New("kata sandi awal minimal 8 karakter")
	}
	return s.repo.CreateStudent(in, in.Password)
}
func (s *Service) CreateLecturer(in LecturerInput) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(in.Email)
	in.EmployeeNumber = strings.TrimSpace(in.EmployeeNumber)
	if in.Name == "" || in.EmployeeNumber == "" || in.StudyProgramID == 0 {
		return errors.New("nama, nomor pegawai, dan program studi wajib diisi")
	}
	if _, e := mail.ParseAddress(in.Email); e != nil {
		return errors.New("alamat email tidak valid")
	}
	if len(in.Password) < 8 {
		return errors.New("kata sandi awal minimal 8 karakter")
	}
	return s.repo.CreateLecturer(in, in.Password)
}
func (s *Service) Toggle(id uint64) error {
	if id == 0 {
		return errors.New("pengguna tidak valid")
	}
	return s.repo.Toggle(id)
}
func studentStatus(v string) bool {
	return v == "active" || v == "leave" || v == "graduated" || v == "inactive"
}
func lecturerStatus(v string) bool { return v == "active" || v == "inactive" || v == "retired" }
