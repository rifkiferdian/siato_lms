package usermanagement

import "testing"

func TestCreateValidatesAdministrativeAccount(t *testing.T) {
	service := NewService(nil)
	tests := []struct {
		name, email, password string
		roleID                uint64
	}{
		{"", "admin@example.com", "Password123!", 1},
		{"Admin", "invalid", "Password123!", 1},
		{"Admin", "admin@example.com", "short", 1},
		{"Admin", "admin@example.com", "Password123!", 0},
	}
	for _, test := range tests {
		if err := service.Create(test.name, test.email, test.password, test.roleID); err == nil {
			t.Fatalf("expected validation error for %#v", test)
		}
	}
}

func TestRoleAndStatusChangesProtectCurrentUser(t *testing.T) {
	service := NewService(nil)
	if err := service.Toggle(7, 7); err == nil {
		t.Fatal("expected self-toggle rejection")
	}
	if err := service.SetRoles(7, 7, []uint64{1}); err == nil {
		t.Fatal("expected self-role rejection")
	}
	if err := service.SetRoles(7, 8, nil); err == nil {
		t.Fatal("expected empty-role rejection")
	}
	if err := service.SetRoles(7, 8, []uint64{1, 1}); err == nil {
		t.Fatal("expected duplicate-role rejection")
	}
}
