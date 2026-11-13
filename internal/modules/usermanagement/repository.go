package usermanagement

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Users(search, status, role string) ([]UserRow, error) {
	q := `SELECT u.id,u.name,u.email,u.is_active,
COALESCE(GROUP_CONCAT(DISTINCT ro.name ORDER BY ro.name SEPARATOR ', '),'Tanpa role') roles,
DATE_FORMAT(u.created_at,'%d/%m/%Y') created_at,
COALESCE(DATE_FORMAT(u.last_login_at,'%d/%m/%Y %H:%i'),'Belum pernah') last_login_at
FROM users u LEFT JOIN model_has_roles mr ON mr.user_id=u.id LEFT JOIN roles ro ON ro.id=mr.role_id WHERE 1=1`
	args := []any{}
	if search != "" {
		like := "%" + search + "%"
		q += " AND (u.name LIKE ? OR u.email LIKE ?)"
		args = append(args, like, like)
	}
	if status == "active" || status == "inactive" {
		q += " AND u.is_active=?"
		args = append(args, status == "active")
	}
	if role != "" {
		q += " AND EXISTS(SELECT 1 FROM model_has_roles mx JOIN roles rx ON rx.id=mx.role_id WHERE mx.user_id=u.id AND rx.name=?)"
		args = append(args, role)
	}
	q += " GROUP BY u.id,u.name,u.email,u.is_active,u.created_at,u.last_login_at ORDER BY u.name"
	var rows []UserRow
	err := r.db.Raw(q, args...).Scan(&rows).Error
	return rows, err
}

func (r *Repository) RoleOptions(adminOnly bool) ([]Option, error) {
	q := "SELECT id,name,CASE name WHEN 'super_admin' THEN 'Super Admin' WHEN 'academic_admin' THEN 'Admin Akademik' WHEN 'lecturer' THEN 'Dosen' WHEN 'student' THEN 'Mahasiswa' ELSE name END label FROM roles"
	if adminOnly {
		q += " WHERE name IN ('super_admin','academic_admin')"
	}
	q += " ORDER BY name"
	var rows []Option
	err := r.db.Raw(q).Scan(&rows).Error
	return rows, err
}

func (r *Repository) Create(name, email, password string, roleID uint64) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var roleName string
		if err := tx.Raw("SELECT name FROM roles WHERE id=? AND name IN ('super_admin','academic_admin')", roleID).Scan(&roleName).Error; err != nil {
			return err
		}
		if roleName == "" {
			return errors.New("role akun administratif tidak valid")
		}
		result := tx.Exec("INSERT INTO users(name,email,password_hash,is_active) VALUES(?,?,?,1)", name, strings.ToLower(email), string(hash))
		if result.Error != nil {
			return result.Error
		}
		var userID uint64
		if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&userID).Error; err != nil {
			return err
		}
		return tx.Exec("INSERT INTO model_has_roles(role_id,user_id) VALUES(?,?)", roleID, userID).Error
	})
}

func (r *Repository) Toggle(userID uint64) error {
	result := r.db.Exec("UPDATE users SET is_active=NOT is_active WHERE id=?", userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("pengguna tidak ditemukan")
	}
	return nil
}

func (r *Repository) SetRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Table("roles").Where("id IN ?", roleIDs).Count(&count).Error; err != nil {
			return err
		}
		if int(count) != len(roleIDs) {
			return errors.New("role tidak valid")
		}
		if err := tx.Exec("DELETE FROM model_has_roles WHERE user_id=?", userID).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if err := tx.Exec("INSERT INTO model_has_roles(role_id,user_id) VALUES(?,?)", roleID, userID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) Roles() ([]RoleRow, error) {
	var rows []RoleRow
	err := r.db.Raw(`SELECT r.id,r.name,CASE r.name WHEN 'super_admin' THEN 'Super Admin' WHEN 'academic_admin' THEN 'Admin Akademik' WHEN 'lecturer' THEN 'Dosen' WHEN 'student' THEN 'Mahasiswa' ELSE r.name END label,
(SELECT COUNT(*) FROM model_has_roles mr WHERE mr.role_id=r.id) user_count,
(SELECT COUNT(*) FROM role_has_permissions rp WHERE rp.role_id=r.id) permission_count FROM roles r ORDER BY r.name`).Scan(&rows).Error
	return rows, err
}

func (r *Repository) Permissions(roleID uint64) ([]PermissionRow, error) {
	var rows []PermissionRow
	err := r.db.Raw(`SELECT p.id,? role_id,p.name,SUBSTRING_INDEX(p.name,'.',1) group_label,
EXISTS(SELECT 1 FROM role_has_permissions rp WHERE rp.role_id=? AND rp.permission_id=p.id) assigned
FROM permissions p ORDER BY group_label,p.name`, roleID, roleID).Scan(&rows).Error
	return rows, err
}

func (r *Repository) SetPermissions(roleID uint64, permissionIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var roleName string
		if err := tx.Raw("SELECT name FROM roles WHERE id=?", roleID).Scan(&roleName).Error; err != nil {
			return err
		}
		if roleName == "" {
			return errors.New("role tidak ditemukan")
		}
		if roleName == "super_admin" {
			return errors.New("permission Super Admin tidak dapat dibatasi")
		}
		if len(permissionIDs) > 0 {
			var count int64
			if err := tx.Table("permissions").Where("id IN ?", permissionIDs).Count(&count).Error; err != nil {
				return err
			}
			if int(count) != len(permissionIDs) {
				return errors.New("permission tidak valid")
			}
		}
		if err := tx.Exec("DELETE FROM role_has_permissions WHERE role_id=?", roleID).Error; err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			if err := tx.Exec("INSERT INTO role_has_permissions(permission_id,role_id) VALUES(?,?)", permissionID, roleID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
