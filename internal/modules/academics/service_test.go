package academics

import "testing"

func TestAcademicYearValidation(t *testing.T) {
	s := NewService(nil)
	for _, in := range []AcademicYearInput{
		{},
		{Name: "2026/2027", StartsOn: "invalid", EndsOn: "2027-06-30"},
		{Name: "2026/2027", StartsOn: "2027-07-01", EndsOn: "2027-06-30"},
	} {
		if err := s.SaveYear(0, in); err == nil {
			t.Fatalf("expected validation error for %+v", in)
		}
	}
}

func TestSemesterValidation(t *testing.T) {
	s := NewService(nil)
	if err := s.SaveSemester(0, SemesterInput{Name: "Ganjil", Code: "2026-1", StartsOn: "2026-08-01", EndsOn: "2026-12-31"}); err == nil {
		t.Fatal("expected missing academic year error")
	}
	if err := s.ActivateSemester(0); err == nil {
		t.Fatal("expected invalid semester error")
	}
}

func TestCourseValidation(t *testing.T) {
	s := NewService(nil)
	for _, credits := range []int{0, 13} {
		if err := s.SaveCourse(0, CourseInput{StudyProgramID: 1, Code: "IF101", Name: "Algoritma", Credits: credits}); err == nil {
			t.Fatalf("expected invalid credits error: %d", credits)
		}
	}
}
