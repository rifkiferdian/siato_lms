package academics

type AcademicYear struct {
	ID                     uint64
	Name, StartsOn, EndsOn string
	IsActive               bool
	SemesterCount          int
}

type Semester struct {
	ID, AcademicYearID                             uint64
	Name, Code, AcademicYearName, StartsOn, EndsOn string
	IsActive                                       bool
	ClassCount                                     int
}

type StudyProgram struct {
	ID                                       uint64
	Code, Name, Degree                       string
	StudentCount, LecturerCount, CourseCount int
}

type Course struct {
	ID, StudyProgramID      uint64
	Code, Name, ProgramName string
	Credits                 int
	ClassCount              int
}

type Option struct {
	ID    uint64
	Label string
}
