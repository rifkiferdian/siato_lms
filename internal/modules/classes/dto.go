package classes

type CreateInput struct {
	CourseID, SemesterID      uint64
	Name, Description, Status string
	Capacity, MaxMeetings     int
}

type PageData struct {
	Portal, Search, Status, Success, Error string
	Classes                                []ClassRow
	Courses, Semesters                     []Option
	CanManage                              bool
}

type DetailData struct {
	Portal, Tab           string
	Class                 ClassDetail
	CanManage             bool
	CanManageParticipants bool
	Students              []StudentRow
	AvailableStudents     []StudentOption
	EntryYears            []uint16
	Success, Error        string
}
