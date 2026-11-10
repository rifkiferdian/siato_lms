package auth

import (
	"gorm.io/gorm"
	"siato_lms/internal/modules/users"
	"strings"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) ByEmail(email string) (users.User, error) {
	var u users.User
	e := r.db.Where("email=?", strings.ToLower(strings.TrimSpace(email))).First(&u).Error
	return u, e
}
func (r *Repository) TouchLogin(id uint64) error {
	return r.db.Exec("UPDATE users SET last_login_at=UTC_TIMESTAMP() WHERE id=?", id).Error
}
