package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"siato_lms/internal/modules/users"
	"strings"
	"time"
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
func (s *Service) Profile(id uint64) (users.User, error) { return s.repo.ByID(id) }
func (s *Service) UpdateProfile(id uint64, name, email string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	if _, e := mail.ParseAddress(email); e != nil {
		return errors.New("alamat email tidak valid")
	}
	return s.repo.UpdateProfile(id, name, email)
}
func validNewPassword(v string) error {
	if len(v) < 8 {
		return errors.New("kata sandi minimal 8 karakter")
	}
	return nil
}
func (s *Service) ChangePassword(id uint64, current, next, confirm string) error {
	if next != confirm {
		return errors.New("konfirmasi kata sandi tidak sama")
	}
	if e := validNewPassword(next); e != nil {
		return e
	}
	u, e := s.repo.ByID(id)
	if e != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(current)) != nil {
		return errors.New("kata sandi saat ini salah")
	}
	hash, e := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if e != nil {
		return e
	}
	return s.repo.UpdatePassword(id, string(hash))
}
func tokenHash(v string) string { sum := sha256.Sum256([]byte(v)); return hex.EncodeToString(sum[:]) }
func (s *Service) RequestReset(email string) (string, error) {
	u, e := s.repo.ByEmail(email)
	if e != nil || !u.IsActive {
		return "", nil
	}
	b := make([]byte, 32)
	if _, e = rand.Read(b); e != nil {
		return "", e
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	e = s.repo.CreateReset(u.ID, tokenHash(token), time.Now().UTC().Add(time.Hour))
	return token, e
}
func (s *Service) ResetPassword(token, password, confirm string) error {
	if token == "" {
		return errors.New("token reset tidak valid")
	}
	if password != confirm {
		return errors.New("konfirmasi kata sandi tidak sama")
	}
	if e := validNewPassword(password); e != nil {
		return e
	}
	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return e
	}
	if e = s.repo.ResetPassword(tokenHash(token), string(hash)); e != nil {
		return errors.New("token reset tidak valid atau kedaluwarsa")
	}
	return nil
}
func (s *Service) Sessions(userID uint64, current string) ([]SessionRow, error) {
	return s.repo.Sessions(userID, current)
}
func (s *Service) RevokeSession(userID, sessionID uint64) error {
	if sessionID == 0 {
		return errors.New("sesi tidak valid")
	}
	return s.repo.RevokeSession(userID, sessionID)
}
