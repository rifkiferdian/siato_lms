package auth

import (
	"gorm.io/gorm"
	"siato_lms/internal/modules/users"
	"strings"
	"time"
)

type SessionRow struct {
	ID                                          uint64
	IPAddress, UserAgent, LastSeenAt, CreatedAt string
	IsCurrent                                   bool
}

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
func (r *Repository) ByID(id uint64) (users.User, error) {
	var u users.User
	e := r.db.First(&u, id).Error
	return u, e
}
func (r *Repository) UpdateProfile(id uint64, name, email string) error {
	return r.db.Exec("UPDATE users SET name=?,email=? WHERE id=?", name, strings.ToLower(email), id).Error
}
func (r *Repository) UpdatePassword(id uint64, hash string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("UPDATE users SET password_hash=? WHERE id=?", hash, id).Error; e != nil {
			return e
		}
		return tx.Exec("UPDATE user_sessions SET revoked_at=UTC_TIMESTAMP() WHERE user_id=? AND revoked_at IS NULL", id).Error
	})
}
func (r *Repository) CreateReset(userID uint64, hash string, expires time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("UPDATE password_reset_tokens SET used_at=UTC_TIMESTAMP() WHERE user_id=? AND used_at IS NULL", userID).Error; e != nil {
			return e
		}
		return tx.Exec("INSERT INTO password_reset_tokens(user_id,token_hash,expires_at) VALUES(?,?,?)", userID, hash, expires).Error
	})
}
func (r *Repository) ResetPassword(tokenHash, passwordHash string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var id, userID uint64
		if e := tx.Raw("SELECT id,user_id FROM password_reset_tokens WHERE token_hash=? AND used_at IS NULL AND expires_at>UTC_TIMESTAMP() FOR UPDATE", tokenHash).Row().Scan(&id, &userID); e != nil {
			return e
		}
		if e := tx.Exec("UPDATE users SET password_hash=? WHERE id=?", passwordHash, userID).Error; e != nil {
			return e
		}
		if e := tx.Exec("UPDATE password_reset_tokens SET used_at=UTC_TIMESTAMP() WHERE id=?", id).Error; e != nil {
			return e
		}
		return tx.Exec("UPDATE user_sessions SET revoked_at=UTC_TIMESTAMP() WHERE user_id=? AND revoked_at IS NULL", userID).Error
	})
}
func (r *Repository) Sessions(userID uint64, currentHash string) ([]SessionRow, error) {
	var rows []SessionRow
	e := r.db.Raw(`SELECT id,COALESCE(ip_address,'-') ip_address,COALESCE(user_agent,'Perangkat tidak dikenal') user_agent,DATE_FORMAT(last_seen_at,'%d/%m/%Y %H:%i') last_seen_at,DATE_FORMAT(created_at,'%d/%m/%Y %H:%i') created_at,(token_hash=?) is_current FROM user_sessions WHERE user_id=? AND revoked_at IS NULL AND expires_at>UTC_TIMESTAMP() ORDER BY is_current DESC,last_seen_at DESC`, currentHash, userID).Scan(&rows).Error
	return rows, e
}
func (r *Repository) RevokeSession(userID, sessionID uint64) error {
	return r.db.Exec("UPDATE user_sessions SET revoked_at=UTC_TIMESTAMP() WHERE id=? AND user_id=?", sessionID, userID).Error
}
