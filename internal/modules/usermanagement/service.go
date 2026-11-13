package usermanagement

import (
	"errors"
	"net/mail"
	"strings"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) Users(search, status, role string) ([]UserRow, []Option, []Option, error) {
	search = strings.TrimSpace(search)
	if status != "active" && status != "inactive" {
		status = ""
	}
	rows, err := s.repo.Users(search, status, strings.TrimSpace(role))
	if err != nil {
		return nil, nil, nil, err
	}
	roles, err := s.repo.RoleOptions(false)
	if err != nil {
		return nil, nil, nil, err
	}
	adminRoles, err := s.repo.RoleOptions(true)
	return rows, roles, adminRoles, err
}
func (s *Service) Create(name, email, password string, roleID uint64) error {
	name, email = strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email))
	if name == "" || len(name) > 120 {
		return errors.New("nama pengguna wajib diisi")
	}
	if _, err := mail.ParseAddress(email); err != nil || len(email) > 190 {
		return errors.New("alamat email tidak valid")
	}
	if len(password) < 8 {
		return errors.New("kata sandi minimal 8 karakter")
	}
	if roleID == 0 {
		return errors.New("role wajib dipilih")
	}
	return s.repo.Create(name, email, password, roleID)
}
func (s *Service) Toggle(actorID, userID uint64) error {
	if userID == 0 {
		return errors.New("pengguna tidak valid")
	}
	if actorID == userID {
		return errors.New("Anda tidak dapat menonaktifkan akun sendiri")
	}
	return s.repo.Toggle(userID)
}
func (s *Service) SetRoles(actorID, userID uint64, roleIDs []uint64) error {
	if userID == 0 || len(roleIDs) == 0 {
		return errors.New("pengguna harus memiliki minimal satu role")
	}
	if actorID == userID {
		return errors.New("Anda tidak dapat mengubah role akun sendiri")
	}
	seen := map[uint64]bool{}
	for _, id := range roleIDs {
		if id == 0 || seen[id] {
			return errors.New("pilihan role tidak valid")
		}
		seen[id] = true
	}
	return s.repo.SetRoles(userID, roleIDs)
}
func (s *Service) Roles(selected uint64) ([]RoleRow, []PermissionRow, RoleRow, error) {
	roles, err := s.repo.Roles()
	if err != nil {
		return nil, nil, RoleRow{}, err
	}
	if selected == 0 && len(roles) > 0 {
		selected = roles[0].ID
	}
	var current RoleRow
	for _, role := range roles {
		if role.ID == selected {
			current = role
			break
		}
	}
	if current.ID == 0 {
		return nil, nil, RoleRow{}, errors.New("role tidak ditemukan")
	}
	permissions, err := s.repo.Permissions(current.ID)
	return roles, permissions, current, err
}
func (s *Service) SetPermissions(roleID uint64, permissionIDs []uint64) error {
	if roleID == 0 {
		return errors.New("role tidak valid")
	}
	seen := map[uint64]bool{}
	for _, id := range permissionIDs {
		if id == 0 || seen[id] {
			return errors.New("pilihan permission tidak valid")
		}
		seen[id] = true
	}
	return s.repo.SetPermissions(roleID, permissionIDs)
}
