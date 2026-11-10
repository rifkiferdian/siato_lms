package rbac

import "gorm.io/gorm"

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

func (s *Service) Roles(userID uint64) ([]string, error) {
	var names []string
	err := s.db.Table("roles r").Select("r.name").Joins("JOIN model_has_roles m ON m.role_id=r.id").Where("m.user_id=?", userID).Order("r.name").Scan(&names).Error
	return names, err
}
func (s *Service) HasPermission(userID uint64, permission string) bool {
	var count int64
	s.db.Raw(`SELECT COUNT(DISTINCT p.id) FROM permissions p
LEFT JOIN role_has_permissions rp ON rp.permission_id=p.id
LEFT JOIN model_has_roles mr ON mr.role_id=rp.role_id AND mr.user_id=?
LEFT JOIN model_has_permissions mp ON mp.permission_id=p.id AND mp.user_id=?
WHERE p.name=? AND (mr.user_id IS NOT NULL OR mp.user_id IS NOT NULL)`, userID, userID, permission).Scan(&count)
	return count > 0
}
