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
