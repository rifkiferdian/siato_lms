package people

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{repo: r} }

const peoplePerPage = 25

func (s *Service) Students(q, status, source string, page int) ([]StudentRow, []ProgramOption, Pagination, PeopleStats, error) {
	if !studentStatus(status) {
		status = ""
	}
	if source != "manual" && source != "siakad" {
		source = ""
	}
	rows, total, page, e := s.repo.Students(strings.TrimSpace(q), status, source, page, peoplePerPage)
	if e != nil {
		return nil, nil, Pagination{}, PeopleStats{}, e
	}
	assignStudentRowNumbers(rows, page)
	programs, e := s.repo.Programs()
	if e != nil {
		return nil, nil, Pagination{}, PeopleStats{}, e
	}
	stats, e := s.repo.StudentStats()
	return rows, programs, newPagination(page, total, peoplePerPage), stats, e
}
func (s *Service) Lecturers(q, status, source string, page int) ([]LecturerRow, []ProgramOption, Pagination, PeopleStats, error) {
	if !lecturerStatus(status) {
		status = ""
	}
	if source != "manual" && source != "siakad" {
		source = ""
	}
	rows, total, page, e := s.repo.Lecturers(strings.TrimSpace(q), status, source, page, peoplePerPage)
	if e != nil {
		return nil, nil, Pagination{}, PeopleStats{}, e
	}
	assignLecturerRowNumbers(rows, page)
	programs, e := s.repo.Programs()
	if e != nil {
		return nil, nil, Pagination{}, PeopleStats{}, e
	}
	stats, e := s.repo.LecturerStats()
	return rows, programs, newPagination(page, total, peoplePerPage), stats, e
}

func newPagination(page, total, perPage int) Pagination {
	totalPages := (total + perPage - 1) / perPage
	if totalPages == 0 {
		totalPages = 1
	}
	p := Pagination{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages, HasPrev: page > 1, HasNext: page < totalPages, PrevPage: page - 1, NextPage: page + 1}
	if total > 0 {
		p.From = (page-1)*perPage + 1
		p.To = p.From + perPage - 1
		if p.To > total {
			p.To = total
		}
	}
	start, end := page-2, page+2
	if start < 1 {
		start = 1
	}
	if end > totalPages {
		end = totalPages
	}
	for i := start; i <= end; i++ {
		p.Pages = append(p.Pages, i)
	}
	return p
}

func assignStudentRowNumbers(rows []StudentRow, page int) {
	for i := range rows {
		rows[i].RowNumber = (page-1)*peoplePerPage + i + 1
	}
}

func assignLecturerRowNumbers(rows []LecturerRow, page int) {
	for i := range rows {
		rows[i].RowNumber = (page-1)*peoplePerPage + i + 1
	}
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
