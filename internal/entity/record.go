package entity

import "time"

type Record struct {
	ID             int `db:"id"`
	UserID         int `db:"user_id"`
	ScheduleItemId int `db:"schedule_item_id"`
	//Body           string    `db:"body"`
	active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
}
