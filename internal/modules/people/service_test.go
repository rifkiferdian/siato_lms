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

func TestPaginationUsesTwentyFiveRowsPerPage(t *testing.T) {
	p := newPagination(2, 61, peoplePerPage)
	if p.PerPage != 25 || p.From != 26 || p.To != 50 || p.TotalPages != 3 {
		t.Fatalf("unexpected pagination: %+v", p)
	}
	if !p.HasPrev || !p.HasNext || p.PrevPage != 1 || p.NextPage != 3 {
		t.Fatalf("unexpected navigation state: %+v", p)
	}
}

func TestRowNumbersContinueAcrossPages(t *testing.T) {
	students := make([]StudentRow, 2)
	lecturers := make([]LecturerRow, 2)
	assignStudentRowNumbers(students, 2)
	assignLecturerRowNumbers(lecturers, 3)
	if students[0].RowNumber != 26 || students[1].RowNumber != 27 {
		t.Fatalf("unexpected student row numbers: %+v", students)
	}
	if lecturers[0].RowNumber != 51 || lecturers[1].RowNumber != 52 {
		t.Fatalf("unexpected lecturer row numbers: %+v", lecturers)
	}
}
