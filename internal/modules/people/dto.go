package people

type StudentInput struct {
	Name, Email, StudentNumber, ExternalID, Password string
	StudyProgramID                                   uint64
	EntryYear                                        int
}
type LecturerInput struct {
	Name, Email, EmployeeNumber, ExternalID, Password string
	StudyProgramID                                    uint64
}
type PageData struct {
	Kind, Search, Status, Source, Success, Error string
	Students                                     []StudentRow
	Lecturers                                    []LecturerRow
	Programs                                     []ProgramOption
	Pagination                                   Pagination
	Stats                                        PeopleStats
}

type PeopleStats struct {
	Total, Active, Inactive, SIAKAD, Manual int
}

type Pagination struct {
	Page, PerPage, Total, TotalPages, From, To, PrevPage, NextPage int
	HasPrev, HasNext                                               bool
	Pages                                                          []int
}
