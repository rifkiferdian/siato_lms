package coursecatalog

import "testing"

func TestCourseComponentsMustEqualTotalCredits(t *testing.T) {
	s := NewService(nil)
	in := CourseInput{CurriculumID: 1, Code: "IF101", Name: "Algoritma", Credits: 3, RecommendedSemester: 1, CourseType: "required", IsActive: true, TheoryCredits: 2}
	if err := s.SaveCourse(0, in); err == nil {
		t.Fatal("expected component total validation error")
	}
}

func TestCourseRejectsInvalidComponentCredits(t *testing.T) {
	s := NewService(nil)
	in := CourseInput{CurriculumID: 1, Code: "IF101", Name: "Algoritma", Credits: 3, RecommendedSemester: 1, CourseType: "required", TheoryCredits: 15, PracticumCredits: -12}
	if err := s.SaveCourse(0, in); err == nil {
		t.Fatal("expected invalid component validation error")
	}
}
