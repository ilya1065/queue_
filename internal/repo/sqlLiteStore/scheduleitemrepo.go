package sqlLiteStore

import (
	"fmt"
	"log/slog"
	"queue/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type ScheduleItemRepo struct {
	db *sqlx.DB
}

//func GetWeek(timeNow time.Time) []entity.ScheduleItem {}

func (repo *ScheduleItemRepo) GetItemByID(id int) ([]entity.ScheduleItem, error) {
	slog.Debug(fmt.Sprintf("GetItemByID: %v", id))
	var item []entity.ScheduleItem
	err := repo.db.Select(&item, `select name, description
 										from schedule_items
 										WHERE id = ? `, id)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	return item, nil
}

func (repo *ScheduleItemRepo) GetItemByTime(start, end time.Time) ([]entity.ScheduleItem, error) {
	slog.Debug("запрос в db GetItemByTime")
	var items []entity.ScheduleItem
	_ = end
	day := start.Format("2006-01-02")
	err := repo.db.Select(&items, `select id, name,description,start_date,end_date
												from schedule_items
												where substr(start_date,1,10) = ?
												  and substr(start_date,12,8) >= '01:00:00'
												order by start_date`, day)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ScheduleItemRepo) GetDay() ([]entity.ScheduleItem, error) {
	slog.Debug("запрос в DB GetDay")
	var items []entity.ScheduleItem
	start, _, err := startEndOfToday()
	if err != nil {
		return nil, err
	}
	day := start.Format("2006-01-02")
	err = repo.db.Select(&items, `select id,name,description,start_date,end_date
												from schedule_items
												where substr(start_date,1,10) = ?
												order by start_date`, day)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func startEndOfToday() (time.Time, time.Time, error) {
	slog.Debug("получение конца и начала дня startEndOfToday")
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	now := time.Now().In(loc)
	y, m, d := now.Date()

	start := time.Date(y, m, d, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)

	return start, end, nil
}
