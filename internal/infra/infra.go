package infra

import (
	"context"
	"log/slog"
	"queue/internal/parser"

	"github.com/jmoiron/sqlx"
)

type Infra struct {
	db *sqlx.DB
}

func NewInfra(db *sqlx.DB) *Infra {
	return &Infra{
		db: db,
	}
}

func (inf *Infra) LoadDBScheduleItem() error {
	slog.Info("загрузка расписания в db")
	url := "https://schedule-of.mirea.ru/_next/data/fR0NO9mu2NSCPRkXv6ZHl/index.json?date=2026-1-16&s=1_4783"
	ev, err := parser.ICSURL(url)
	if err != nil {
		return err
	}
	ctx := context.Background()
	tx, err := inf.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO schedule_items
    											(name, description,start_date,end_date,external_id)
												values(?,?,?,?,?)
												on conflict (external_id,start_date) DO update set
												name=excluded.name,
												description=excluded.description,
												start_date=excluded.start_date,
												end_date=excluded.end_date`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, e := range ev {
		_, err = stmt.ExecContext(ctx, e.Name, e.Description, e.StartDate, e.EndDate, e.ExternalID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
