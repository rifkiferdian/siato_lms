package usermanagement

type UserRow struct {
	ID                                         uint64
	Name, Email, Roles, CreatedAt, LastLoginAt string
	IsActive                                   bool
}

type RoleRow struct {
	ID                         uint64
	Name, Label                string
	UserCount, PermissionCount int
}

type PermissionRow struct {
	ID, RoleID       uint64
	Name, GroupLabel string
	Assigned         bool
}

type Option struct {
	ID    uint64
	Name  string
	Label string
}

type UsersPage struct {
	Search, Status, Role, Success, Error string
	Users                                []UserRow
	Roles, AdminRoles                    []Option
}

type RolesPage struct {
	Roles          []RoleRow
	Permissions    []PermissionRow
	SelectedRole   RoleRow
	Success, Error string
}
