package schedules

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

var ErrConflict = errors.New("kelas atau ruangan memiliki jadwal yang bertabrakan pada waktu tersebut")

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }
func (s *Service) Page(portal string, userID uint64, day int) ([]Row, []ClassOption, error) {
	rows, e := s.repo.List(portal, userID, day)
	if e != nil {
		return nil, nil, e
	}
	classes, e := s.repo.Classes()
	return rows, classes, e
}
func (s *Service) Create(in CreateInput) error {
	in.Room = strings.TrimSpace(in.Room)
	in.DeliveryMode = strings.ToLower(in.DeliveryMode)
	if in.ClassID == 0 || in.DayOfWeek < 1 || in.DayOfWeek > 7 || in.Room == "" {
		return errors.New("kelas, hari, dan ruangan wajib diisi")
	}
	start, e1 := time.Parse("15:04", in.StartTime)
	end, e2 := time.Parse("15:04", in.EndTime)
	if e1 != nil || e2 != nil || !end.After(start) {
		return errors.New("rentang waktu jadwal tidak valid")
	}
	if in.DeliveryMode != "onsite" && in.DeliveryMode != "online" && in.DeliveryMode != "hybrid" {
		return errors.New("mode perkuliahan tidak valid")
	}
	if in.MeetingURL != "" {
		u, e := url.ParseRequestURI(in.MeetingURL)
		if e != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return errors.New("tautan pertemuan tidak valid")
		}
	}
	conflict, e := s.repo.HasConflict(in)
	if e != nil {
		return e
	}
	if conflict {
		return ErrConflict
	}
	return s.repo.Create(in)
}
func (s *Service) Delete(id uint64) error {
	if id == 0 {
		return errors.New("jadwal tidak valid")
	}
	return s.repo.Delete(id)
}
