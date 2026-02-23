package sqlLiteStore

import (
	"fmt"
	"log/slog"
	"queue/internal/config"
	"queue/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type ScheduleItemRepo struct {
	db  *sqlx.DB
	cfg *config.Config
}

//func GetWeek(timeNow time.Time) []entity.ScheduleItem {}

// GetItemByID возвращает пару по id
func (repo *ScheduleItemRepo) GetItemByID(id int) (*entity.ScheduleItem, error) {
	slog.Debug(fmt.Sprintf("GetItemByID: %v", id))
	var item entity.ScheduleItem
	err := repo.db.Get(&item, `select name, description
 										from schedule_items
 										WHERE id = ? `, id)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	return &item, nil
}

// GetItemByTime возвращает срез пар на день даты start
func (repo *ScheduleItemRepo) GetItemByTime(start, end time.Time) ([]entity.ScheduleItem, error) {
	slog.Debug("запрос в db GetItemByTime")
	var items []entity.ScheduleItem
	_ = end
	fmt.Println(start)
	day := start.Format("2006-01-02")
	fmt.Println(day)
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

func (repo *ScheduleItemRepo) UpdateScheduleForTime(items []entity.ScheduleItem, start time.Time) error {
	slog.Debug(fmt.Sprintf("UpdateScheduleForTime: %v", start))

	tx, err := repo.db.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("UpdateScheduleForTime: %v", err))
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`delete from schedule_items
									WHERE start_date >= ?`, start)
	if err != nil {
		slog.Error(fmt.Sprintf("UpdateScheduleForTime: %v", err))
		tx.Rollback()
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO schedule_items (name, description, start_date, end_date, external_id)
									VALUES (?,?,?,?,?)
									ON CONFLICT(name, start_date, end_date)
									DO UPDATE SET
  									description = excluded.description,
  									external_id = excluded.external_id;`)
	if err != nil {
		slog.Error(fmt.Sprintf("ошибка подготовки запроса UpdateScheduleForTime: %v", err))
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		_, err = stmt.Exec(item.Name, item.Description, item.StartDate, item.EndDate, item.ExternalID)
		if err != nil {
			slog.Debug("ошибка при выполнении загрузки расписания UpdateScheduleForTime")
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		slog.Debug(fmt.Sprintf("UpdateScheduleForTime: %v", err))
		return err
	}

	return nil
}

/*
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

// Возвращает начало и конец дня
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
*/
