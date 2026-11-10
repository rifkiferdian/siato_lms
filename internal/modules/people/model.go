package people

type StudentRow struct {
	ID, UserID                                                                                                 uint64
	Name, Email, StudentNumber, ProgramCode, ProgramName, AcademicStatus, DataSource, ExternalID, LastSyncedAt string
	EntryYear                                                                                                  uint16
	IsActive                                                                                                   bool
	ClassCount                                                                                                 int
}
type LecturerRow struct {
	ID, UserID                                                                                                  uint64
	Name, Email, EmployeeNumber, ProgramCode, ProgramName, AcademicStatus, DataSource, ExternalID, LastSyncedAt string
	IsActive                                                                                                    bool
	ClassCount                                                                                                  int
}
type ProgramOption struct {
	ID    uint64
	Label string
}
