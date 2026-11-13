package classes

import "testing"

func TestCreateRejectsInvalidClassBeforeRepository(t *testing.T) {
	s := NewService(nil)
	tests := []CreateInput{
		{SemesterID: 1, Name: "A", Capacity: 40, MaxMeetings: 16, Status: "active"},
		{CourseID: 1, SemesterID: 1, Name: "", Capacity: 40, MaxMeetings: 16, Status: "active"},
		{CourseID: 1, SemesterID: 1, Name: "A", Capacity: 0, MaxMeetings: 16, Status: "active"},
		{CourseID: 1, SemesterID: 1, Name: "A", Capacity: 40, MaxMeetings: 40, Status: "active"},
		{CourseID: 1, SemesterID: 1, Name: "A", Capacity: 40, MaxMeetings: 16, Status: "unknown"},
	}
	for i, input := range tests {
		if err := s.Create(input); err == nil {
			t.Fatalf("case %d: expected validation error", i)
		}
	}
}

func TestAddStudentsRejectsInvalidInputBeforeRepository(t *testing.T) {
	s := NewService(nil)
	for _, ids := range [][]uint64{nil, {}, {0}, {1, 1}} {
		if err := s.AddStudents(1, ids); err == nil {
			t.Fatalf("expected validation error for IDs %v", ids)
		}
	}
	if err := s.AddStudents(0, []uint64{1}); err == nil {
		t.Fatal("expected validation error for empty class ID")
	}
}

func TestClassStudentManagers(t *testing.T) {
	if !canManageClassStudents([]string{"academic_admin"}) {
		t.Fatal("academic admin should manage class students")
	}
	if !canManageClassStudents([]string{"super_admin"}) {
		t.Fatal("super admin should manage class students")
	}
	if canManageClassStudents([]string{"lecturer"}) || canManageClassStudents([]string{"student"}) {
		t.Fatal("lecturer and student must not manage class students")
	}
}
