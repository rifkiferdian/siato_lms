package coursecatalog

type Course struct {
	RowNumber                                                      int `gorm:"-"`
	ID, CurriculumID                                               uint64
	Code, Name, ProgramName, CurriculumName, CourseType            string
	Credits, RecommendedSemester                                   int
	TheoryCredits, PracticumCredits, PracticeCredits, FieldCredits int
	IsActive                                                       bool
}
type Curriculum struct {
	ID, StudyProgramID      uint64
	Code, Name, ProgramName string
	EffectiveYear           int
	IsActive                bool
}
type Option struct {
	ID    uint64
	Label string
}
type PageData struct {
	Search                                      string
	ProgramID, CurriculumID                     uint64
	Page, PrevPage, NextPage, TotalPages, Total int
	Courses                                     []Course
	Curricula                                   []Curriculum
	Programs                                    []Option
	Success, Error                              string
	CanManage                                   bool
}
type CourseInput struct {
	CurriculumID                                                   uint64
	Code, Name, CourseType                                         string
	Credits, RecommendedSemester                                   int
	TheoryCredits, PracticumCredits, PracticeCredits, FieldCredits int
	IsActive                                                       bool
}
type CurriculumInput struct {
	StudyProgramID uint64
	Code, Name     string
	EffectiveYear  int
}
