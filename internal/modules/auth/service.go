package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"siato_lms/internal/modules/users"
)

var ErrInvalidCredentials = errors.New("email atau kata sandi salah")

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }
func (s *Service) Login(email, password string) (users.User, error) {
	u, e := s.repo.ByEmail(email)
	if e != nil || !u.IsActive || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return users.User{}, ErrInvalidCredentials
	}
	_ = s.repo.TouchLogin(u.ID)
	return u, nil
}
