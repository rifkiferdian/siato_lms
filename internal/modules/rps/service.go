package rps

import (
	"errors"
	"strings"
)

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{r} }
func (s *Service) Page(courseID, rpsID, userID uint64, portal string) (PageData, error) {
	info, manage, e := s.repo.Course(courseID, userID, portal)
	if e != nil {
		return PageData{}, e
	}
	versions, e := s.repo.Versions(courseID)
	if e != nil {
		return PageData{}, e
	}
	doc, rows, e := s.repo.Get(courseID, rpsID)
	if portal == "student" && doc.Status != "active" {
		doc = RPS{}
		rows = nil
	}
	return PageData{Portal: portal, Course: info, RPS: doc, Versions: versions, CPMKs: rows, CanManage: manage}, e
}
func clean(v string) string { return strings.TrimSpace(v) }
func (s *Service) Save(courseID, rpsID, userID uint64, in RPSInput) error {
	in.Version = clean(in.Version)
	in.Description = clean(in.Description)
	in.LearningMethods = clean(in.LearningMethods)
	in.MainReferences = clean(in.MainReferences)
	in.SupportingReferences = clean(in.SupportingReferences)
	if courseID == 0 || in.Version == "" || in.Description == "" {
		return errors.New("versi dan deskripsi RPS wajib diisi")
	}
	return s.repo.Save(courseID, rpsID, userID, in)
}
func (s *Service) Activate(courseID, rpsID, userID uint64) error {
	doc, rows, e := s.repo.Get(courseID, rpsID)
	if e != nil {
		return e
	}
	if doc.ID == 0 {
		return errors.New("RPS belum dibuat")
	}
	if len(rows) == 0 {
		return errors.New("tambahkan minimal satu CPMK sebelum mengaktifkan RPS")
	}
	return s.repo.Activate(courseID, rpsID, userID)
}
func (s *Service) SaveCPMK(courseID, rpsID, id uint64, in CPMKInput) error {
	in.Code = strings.ToUpper(clean(in.Code))
	in.Title = clean(in.Title)
	in.Description = clean(in.Description)
	if in.Code == "" || in.Title == "" || in.Description == "" {
		return errors.New("kode, judul, dan deskripsi CPMK wajib diisi")
	}
	if in.SortOrder < 1 || in.SortOrder > 999 {
		return errors.New("urutan CPMK tidak valid")
	}
	doc, _, e := s.repo.Get(courseID, rpsID)
	if e != nil || doc.ID == 0 {
		return errors.New("simpan RPS terlebih dahulu")
	}
	return s.repo.SaveCPMK(courseID, rpsID, id, in)
}
func (s *Service) DeleteCPMK(courseID, rpsID, id uint64) error {
	if courseID == 0 || rpsID == 0 || id == 0 {
		return errors.New("CPMK tidak valid")
	}
	return s.repo.DeleteCPMK(courseID, rpsID, id)
}
