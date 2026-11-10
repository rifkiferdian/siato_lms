package schedules

type CreateInput struct {
	ClassID                                                                           uint64
	DayOfWeek                                                                         int
	StartTime, EndTime, Room, DeliveryMode, MeetingURL, EffectiveFrom, EffectiveUntil string
}

type PageData struct {
	Portal         string
	Schedules      []Row
	Classes        []ClassOption
	SelectedDay    int
	Days           []DayOption
	CanManage      bool
	Success, Error string
}

type DayOption struct {
	Value int
	Label string
}

func Days() []DayOption {
	return []DayOption{{1, "Senin"}, {2, "Selasa"}, {3, "Rabu"}, {4, "Kamis"}, {5, "Jumat"}, {6, "Sabtu"}, {7, "Minggu"}}
}
