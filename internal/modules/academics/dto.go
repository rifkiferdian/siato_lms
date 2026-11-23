package academics

type AcademicYearInput struct{ Name, StartsOn, EndsOn string }
type SemesterInput struct {
	AcademicYearID               uint64
	Name, Code, StartsOn, EndsOn string
}
type StudyProgramInput struct{ Code, Name, Degree string }
type CourseInput struct {
	StudyProgramID uint64
	Code, Name     string
	Credits        int
}

type PageData struct {
	Tab, Success, Error         string
	CanManage                   bool
	AcademicYears               []AcademicYear
	Semesters                   []Semester
	StudyPrograms               []StudyProgram
	Courses                     []Course
	YearOptions, ProgramOptions []Option
}
