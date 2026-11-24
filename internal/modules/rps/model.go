package rps

type RPS struct {
	ID, CourseID                                                                                                uint64
	Version, Description, LearningMethods, MainReferences, SupportingReferences, Status, ActivatedAt, UpdatedAt string
}
type CPMK struct {
	ID                       uint64
	Code, Title, Description string
	SortOrder                int
}
type CourseInfo struct {
	ID                                                                                           uint64
	Code, Name, ProgramName, CurriculumName, CourseType                                          string
	Credits, RecommendedSemester, TheoryCredits, PracticumCredits, PracticeCredits, FieldCredits int
	IsActive                                                                                     bool
}
type PageData struct {
	Portal         string
	Course         CourseInfo
	RPS            RPS
	Versions       []RPS
	CPMKs          []CPMK
	CanManage      bool
	Success, Error string
}
type RPSInput struct{ Version, Description, LearningMethods, MainReferences, SupportingReferences string }
type CPMKInput struct {
	Code, Title, Description string
	SortOrder                int
}
