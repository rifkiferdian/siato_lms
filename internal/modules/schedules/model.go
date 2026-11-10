package schedules

import "time"

type Schedule struct {
	ID             uint64 `gorm:"primaryKey"`
	ClassID        uint64
	DayOfWeek      uint8
	StartTime      string
	EndTime        string
	Room           string
	DeliveryMode   string
	MeetingURL     *string
	EffectiveFrom  *time.Time
	EffectiveUntil *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Schedule) TableName() string { return "class_schedules" }

type Row struct {
	ID                                                                               uint64
	ClassID                                                                          uint64
	CourseCode, CourseName, ClassName, LecturerNames, Room, DeliveryMode, MeetingURL string
	DayOfWeek                                                                        uint8
	DayName, StartTime, EndTime                                                      string
	Credits                                                                          uint8
}

type ClassOption struct {
	ID    uint64
	Label string
}
