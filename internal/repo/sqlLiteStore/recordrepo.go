package sqlLiteStore

import (
	"log/slog"
	"queue/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type RecordRepo struct {
	db *sqlx.DB
}

func (repo *RecordRepo) DeleteRecord(userID int64, scheduleID int) error {
	slog.Debug(`RecordRepo.DeleteRecord()`)

	res, err := repo.db.Exec(`DELETE from records
								where user_id = $1 and schedule_item_id =$2 `, userID, scheduleID)
	if err != nil {
		slog.Error(`RecordRepo.DeleteRecord()`, err)
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return entity.ErrAlreadyRegistered
	}

	return nil
}

func (repo RecordRepo) AddUserToItem(id int64, scheduleItemID int) error {
	slog.Debug("работа с db RecordRepo.AddUserToItem")
	timeNow := time.Now()
	res, err := repo.db.Exec(`insert into records (user_id,schedule_item_id,createdAt)
						 		  values (?, ?,?)
						 		   ON CONFLICT(user_id, schedule_item_id) DO NOTHING`, id, scheduleItemID, timeNow)
	if err != nil {
		slog.Debug(err.Error())
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return entity.ErrAlreadyRegistered
	}

	return nil
}

func (repo RecordRepo) GetUserByItemID(id int) ([]entity.User, error) {
	slog.Debug("работа с db RecordRepo.GetUserByItemID ")
	var users []entity.User
	err := repo.db.Select(&users, `select u.id,u.name
										 from records r
										 join users u ON u.id = r.user_id
										 where r.schedule_item_id = ? 
										 order by r.createdAt asc`, id)
	if err != nil {
		slog.Warn(err.Error())
		return nil, err
	}
	return users, nil
}
