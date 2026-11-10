package schedules

import "testing"

func TestCreateRejectsInvalidScheduleBeforeRepository(t *testing.T) {
	s := NewService(nil)
	tests := []CreateInput{
		{ClassID: 1, DayOfWeek: 0, StartTime: "08:00", EndTime: "10:00", Room: "A", DeliveryMode: "onsite"},
		{ClassID: 1, DayOfWeek: 1, StartTime: "10:00", EndTime: "08:00", Room: "A", DeliveryMode: "onsite"},
		{ClassID: 1, DayOfWeek: 1, StartTime: "08:00", EndTime: "10:00", Room: "A", DeliveryMode: "invalid"},
		{ClassID: 1, DayOfWeek: 1, StartTime: "08:00", EndTime: "10:00", Room: "A", DeliveryMode: "online", MeetingURL: "javascript:alert(1)"},
	}
	for i, input := range tests {
		if err := s.Create(input); err == nil {
			t.Fatalf("case %d: expected validation error", i)
		}
	}
}
