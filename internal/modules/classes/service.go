package classes

import (
	"errors"
	"strings"
)

var ErrClassHasStudents = errors.New("kelas yang memiliki peserta tidak dapat dihapus")

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{repo: r} }
func (s *Service) List(portal string, userID uint64, search, status string) ([]ClassRow, []Option, []Option, error) {
	if status != "" && !validStatus(status) {
		status = ""
	}
	rows, e := s.repo.List(portal, userID, strings.TrimSpace(search), status)
	if e != nil {
		return nil, nil, nil, e
	}
	courses, semesters, e := s.repo.Options()
	return rows, courses, semesters, e
}
func (s *Service) Detail(id, userID uint64, portal string) (ClassDetail, error) {
	if id == 0 {
		return ClassDetail{}, gormErrNotFound
	}
	return s.repo.Detail(id, userID, portal)
}
func (s *Service) Create(in CreateInput) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	if in.CourseID == 0 || in.SemesterID == 0 || in.Name == "" {
		return errors.New("mata kuliah, semester, dan nama kelas wajib diisi")
	}
	if in.Capacity < 1 || in.Capacity > 500 {
		return errors.New("kapasitas kelas harus antara 1 sampai 500")
	}
	if in.MaxMeetings < 1 || in.MaxMeetings > 32 {
		return errors.New("jumlah pertemuan harus antara 1 sampai 32")
	}
	if !validStatus(in.Status) {
		return errors.New("status kelas tidak valid")
	}
	return s.repo.Create(in)
}
func (s *Service) Delete(id uint64) error {
	if id == 0 {
		return errors.New("kelas tidak valid")
	}
	return s.repo.Delete(id)
}
func validStatus(v string) bool {
	return v == "draft" || v == "active" || v == "completed" || v == "cancelled"
}

var gormErrNotFound = errors.New("kelas tidak ditemukan atau tidak dapat diakses")
