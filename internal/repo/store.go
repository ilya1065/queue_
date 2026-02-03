package repo

type Store interface {
	ScheduleItem() ScheduleRepo
	User() UserRepo
	Record() RecordsRepo
}
