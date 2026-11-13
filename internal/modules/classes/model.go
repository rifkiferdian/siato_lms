package classes

type ClassRow struct {
	ID                                                                                                                 uint64
	CourseCode, CourseName, ClassName, Description, Status, SemesterName, AcademicYear, LecturerNames, ScheduleSummary string
	Credits, MaxMeetings                                                                                               uint8
	Capacity, StudentCount, LecturerCount, ScheduleCount                                                               int
}

type ClassDetail struct {
	ClassRow
	ProgramName string
}

type Option struct {
	ID    uint64
	Label string
}

type StudentRow struct {
	ID                                      uint64
	Name, Email, StudentNumber, ProgramName string
}
