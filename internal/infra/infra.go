package infra

import (
	"context"
	"log/slog"
	"queue/internal/config"
	"queue/internal/parser"
	"time"

	"github.com/jmoiron/sqlx"
)

type Infra struct {
	db  *sqlx.DB
	cfg *config.Config
}

func NewInfra(db *sqlx.DB, config *config.Config) *Infra {
	return &Infra{
		db:  db,
		cfg: config,
	}
}

func (inf *Infra) LoadDBScheduleItem(start, end time.Time) error {
	slog.Info("загрузка расписания в db")
	url := inf.cfg.SchedulerURL
	ev, err := parser.ICSURL(start, end, url)
	if err != nil {
		return err
	}
	ctx := context.Background()
	tx, err := inf.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//compareItem(tx)\
	// подготовка запроса
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO schedule_items
    											(name, description,start_date,end_date,external_id)
												values(?,?,?,?,?)
												on conflict (name,start_date,end_date) DO update set
												name=excluded.name,
												description=excluded.description,
												start_date=excluded.start_date,
												end_date=excluded.end_date,
												external_id=excluded.external_id`)
	//updateStmt, err := tx.PrepareContext(ctx, `UPDATE schedule_items SET name=excluded.name,`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, e := range ev {
		// реализация подготовленного запроса
		_, err = stmt.ExecContext(ctx, e.Name, e.Description, e.StartDate, e.EndDate, e.ExternalID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

/*// Есть пара в это время 1
// нет пары в это время 0
func compareItem(tx *sqlx.Tx, start, end *time.Time) (bool, error) {
	var exists bool
	err := tx.Select(&exists, `select exists(select 1 from schedule_items where start_date = ? and end_date = ?  )`, start, end)
	if err != nil {
		return false, err
	}
	return exists, nil
}
*/
