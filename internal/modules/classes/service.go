package classes

import (
	"errors"
	"math"
	"strings"
)

var ErrClassHasStudents = errors.New("kelas yang memiliki peserta tidak dapat dihapus")
var ErrClassFull = errors.New("kapasitas kelas tidak mencukupi")
var ErrInvalidStudent = errors.New("mahasiswa tidak valid atau akun tidak aktif")
var ErrStudentAlreadyEnrolled = errors.New("mahasiswa sudah terdaftar di kelas")
var ErrEnrollmentNotFound = errors.New("peserta kelas tidak ditemukan")

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
func (s *Service) Students(classID uint64) ([]StudentRow, []StudentOption, error) {
	if classID == 0 {
		return nil, nil, gormErrNotFound
	}
	return s.repo.Students(classID)
}
func (s *Service) AddStudents(classID uint64, studentIDs []uint64) error {
	if classID == 0 || len(studentIDs) == 0 {
		return errors.New("pilih minimal satu mahasiswa")
	}
	seen := make(map[uint64]bool, len(studentIDs))
	for _, id := range studentIDs {
		if id == 0 || seen[id] {
			return ErrInvalidStudent
		}
		seen[id] = true
	}
	return s.repo.AddStudents(classID, studentIDs)
}
func (s *Service) RemoveStudent(classID, studentID uint64) error {
	if classID == 0 || studentID == 0 {
		return ErrEnrollmentNotFound
	}
	return s.repo.RemoveStudent(classID, studentID)
}
func (s *Service) Create(in CreateInput) error {
	if err := validateClass(&in); err != nil {
		return err
	}
	return s.repo.Create(in)
}
func validateClass(in *CreateInput) error {
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
	return nil
}
func (s *Service) Update(id uint64, in CreateInput) error {
	if id == 0 {
		return errors.New("kelas tidak valid")
	}
	if err := validateClass(&in); err != nil {
		return err
	}
	return s.repo.Update(id, in)
}
func (s *Service) Lecturers(classID uint64) ([]LecturerOption, error) {
	if classID == 0 {
		return nil, gormErrNotFound
	}
	return s.repo.Lecturers(classID)
}
func (s *Service) SetLecturers(classID uint64, assignments []LecturerAssignment, coordinator uint64) error {
	if classID == 0 || len(assignments) == 0 {
		return errors.New("pilih minimal satu dosen pengampu")
	}
	seen := map[uint64]bool{}
	found := false
	validComponents := map[string]bool{"T": true, "P": true, "PR": true, "L": true}
	required, err := s.repo.ComponentCredits(classID)
	if err != nil {
		return err
	}
	allocated := map[string]float64{}
	for _, assignment := range assignments {
		id := assignment.LecturerID
		if id == 0 || seen[id] {
			return errors.New("pilihan dosen tidak valid")
		}
		if len(assignment.Components) == 0 {
			return errors.New("setiap dosen pengampu wajib memiliki minimal satu komponen mengajar")
		}
		componentSeen := map[string]bool{}
		for _, component := range assignment.Components {
			if !validComponents[component.Code] || componentSeen[component.Code] || component.Credits <= 0 || component.Credits > 12 {
				return errors.New("komponen mengajar tidak valid")
			}
			componentSeen[component.Code] = true
			allocated[component.Code] += component.Credits
		}
		seen[id] = true
		if id == coordinator {
			found = true
		}
	}
	if !found {
		return errors.New("koordinator harus dipilih dari dosen pengampu")
	}
	for code, credits := range required {
		if math.Abs(allocated[code]-credits) > 0.001 {
			return errors.New("total alokasi SKS dosen untuk komponen " + code + " harus sama dengan SKS komponen mata kuliah")
		}
	}
	for code, credits := range allocated {
		if _, ok := required[code]; !ok && credits > 0 {
			return errors.New("komponen " + code + " tidak tersedia pada mata kuliah")
		}
	}
	return s.repo.SetLecturers(classID, assignments, coordinator)
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
