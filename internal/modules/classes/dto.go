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
	Portal    string
	Class     ClassDetail
	CanManage bool
}
