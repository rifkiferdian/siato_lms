package people

import "testing"

func TestCreateStudentValidatesInputBeforeRepository(t *testing.T) {
	s := NewService(nil)
	input := StudentInput{Name: "Mahasiswa", Email: "bukan-email", StudentNumber: "20260001", StudyProgramID: 1, EntryYear: 2026, Password: "Password123!"}
	if err := s.CreateStudent(input); err == nil {
		t.Fatal("expected invalid email to be rejected")
	}
}

func TestCreateLecturerRequiresStrongEnoughInitialPassword(t *testing.T) {
	s := NewService(nil)
	input := LecturerInput{Name: "Dosen", Email: "dosen@example.test", EmployeeNumber: "DSN001", StudyProgramID: 1, Password: "short"}
	if err := s.CreateLecturer(input); err == nil {
		t.Fatal("expected short password to be rejected")
	}
}
