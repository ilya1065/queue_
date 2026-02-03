package entity

import "time"

type ScheduleItem struct {
	Id          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	StartDate   *time.Time `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
	ExternalID  string     `db:"external_id"`
}
