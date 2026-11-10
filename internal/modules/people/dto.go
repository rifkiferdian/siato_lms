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
}
