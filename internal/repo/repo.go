package repo

import (
	"queue/internal/entity"
	"time"
)

type ScheduleRepo interface {
	GetDay() ([]entity.ScheduleItem, error)
	GetItemByTime(start, end time.Time) ([]entity.ScheduleItem, error)
	GetItemByID(id int) ([]entity.ScheduleItem, error)
}

type UserRepo interface {
	NewUser(id int64, name string) error
	RenameUser(id int64, newName string) error
	Exists(id int64) (bool, error)
}

type RecordsRepo interface {
	AddUserToItem(id int64, scheduleItemID int) error
	GetUserByItemID(id int) ([]entity.User, error)
}
