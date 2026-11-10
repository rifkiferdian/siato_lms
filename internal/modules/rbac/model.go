package rbac

type Role struct {
	ID        uint64
	Name      string
	GuardName string
}

func (Role) TableName() string { return "roles" }

type Permission struct {
	ID        uint64
	Name      string
	GuardName string
}

func (Permission) TableName() string { return "permissions" }
