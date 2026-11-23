package academics

import (
	"errors"
	"strings"
	"time"
)

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{repo: r} }
func (s *Service) Page() (PageData, error) {
	y, se, p, c, e := s.repo.List()
	d := PageData{AcademicYears: y, Semesters: se, StudyPrograms: p, Courses: c}
	for _, v := range y {
		d.YearOptions = append(d.YearOptions, Option{v.ID, v.Name})
	}
	for _, v := range p {
		d.ProgramOptions = append(d.ProgramOptions, Option{v.ID, v.Code + " - " + v.Name})
	}
	return d, e
}
func clean(v string) string { return strings.TrimSpace(v) }
func validDates(start, end string) bool {
	a, e1 := time.Parse("2006-01-02", start)
	b, e2 := time.Parse("2006-01-02", end)
	return e1 == nil && e2 == nil && !b.Before(a)
}
func validateYear(in *AcademicYearInput) error {
	in.Name = clean(in.Name)
	if in.Name == "" || !validDates(in.StartsOn, in.EndsOn) {
		return errors.New("nama dan rentang tanggal tahun akademik tidak valid")
	}
	return nil
}
func validateSemester(in *SemesterInput) error {
	in.Name = clean(in.Name)
	in.Code = clean(in.Code)
	if in.AcademicYearID == 0 || in.Name == "" || in.Code == "" || !validDates(in.StartsOn, in.EndsOn) {
		return errors.New("tahun akademik, nama, kode, dan rentang tanggal semester wajib valid")
	}
	return nil
}
func validateProgram(in *StudyProgramInput) error {
	in.Code = clean(in.Code)
	in.Name = clean(in.Name)
	in.Degree = clean(in.Degree)
	if in.Code == "" || in.Name == "" || in.Degree == "" {
		return errors.New("kode, nama, dan jenjang program studi wajib diisi")
	}
	return nil
}
func validateCourse(in *CourseInput) error {
	in.Code = clean(in.Code)
	in.Name = clean(in.Name)
	if in.StudyProgramID == 0 || in.Code == "" || in.Name == "" || in.Credits < 1 || in.Credits > 12 {
		return errors.New("program studi, kode, nama, dan SKS 1–12 wajib diisi")
	}
	return nil
}
func (s *Service) SaveYear(id uint64, in AcademicYearInput) error {
	if e := validateYear(&in); e != nil {
		return e
	}
	if id == 0 {
		return s.repo.CreateYear(in)
	}
	return s.repo.UpdateYear(id, in)
}
func (s *Service) DeleteYear(id uint64) error { return s.repo.DeleteYear(id) }
func (s *Service) SaveSemester(id uint64, in SemesterInput) error {
	if e := validateSemester(&in); e != nil {
		return e
	}
	if id == 0 {
		return s.repo.CreateSemester(in)
	}
	return s.repo.UpdateSemester(id, in)
}
func (s *Service) DeleteSemester(id uint64) error { return s.repo.DeleteSemester(id) }
func (s *Service) ActivateSemester(id uint64) error {
	if id == 0 {
		return errors.New("semester tidak valid")
	}
	return s.repo.ActivateSemester(id)
}
func (s *Service) SaveProgram(id uint64, in StudyProgramInput) error {
	if e := validateProgram(&in); e != nil {
		return e
	}
	if id == 0 {
		return s.repo.CreateProgram(in)
	}
	return s.repo.UpdateProgram(id, in)
}
func (s *Service) DeleteProgram(id uint64) error { return s.repo.DeleteProgram(id) }
func (s *Service) SaveCourse(id uint64, in CourseInput) error {
	if e := validateCourse(&in); e != nil {
		return e
	}
	if id == 0 {
		return s.repo.CreateCourse(in)
	}
	return s.repo.UpdateCourse(id, in)
}
func (s *Service) DeleteCourse(id uint64) error { return s.repo.DeleteCourse(id) }
