package coursecatalog

import (
	"errors"
	"strings"
	"time"
)

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{r} }
func (s *Service) Page(q string, p, k uint64, page int) (PageData, error) {
	return s.repo.Page(strings.TrimSpace(q), p, k, page)
}
func (s *Service) SaveCourse(id uint64, in CourseInput) error {
	in.Code = strings.TrimSpace(in.Code)
	in.Name = strings.TrimSpace(in.Name)
	if in.CurriculumID == 0 || in.Code == "" || in.Name == "" || in.Credits < 1 || in.Credits > 12 || in.RecommendedSemester < 1 || in.RecommendedSemester > 14 || (in.CourseType != "required" && in.CourseType != "elective") {
		return errors.New("data mata kuliah tidak valid")
	}
	components := []int{in.TheoryCredits, in.PracticumCredits, in.PracticeCredits, in.FieldCredits}
	total := 0
	for _, credits := range components {
		if credits < 0 || credits > 12 {
			return errors.New("SKS komponen tidak valid")
		}
		total += credits
	}
	if total != in.Credits {
		return errors.New("jumlah SKS Teori, Praktikum, Praktik, dan Lapangan harus sama dengan total SKS")
	}
	return s.repo.SaveCourse(id, in)
}
func (s *Service) DeleteCourse(id uint64) error {
	if id == 0 {
		return errors.New("mata kuliah tidak valid")
	}
	return s.repo.DeleteCourse(id)
}
func (s *Service) SaveCurriculum(id uint64, in CurriculumInput) error {
	in.Code = strings.TrimSpace(in.Code)
	in.Name = strings.TrimSpace(in.Name)
	if in.StudyProgramID == 0 || in.Code == "" || in.Name == "" || in.EffectiveYear < 2000 || in.EffectiveYear > time.Now().Year()+5 {
		return errors.New("data kurikulum tidak valid")
	}
	return s.repo.SaveCurriculum(id, in)
}
