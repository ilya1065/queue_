package app

import (
	"log/slog"
	"queue/internal/config"
	"queue/internal/infra"
	"queue/internal/repo/sqlLiteStore"
	"queue/internal/server"
	"queue/internal/tgbot"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Start(cfg *config.Config) error {
	slog.Info("Старт сервиса")
	db, err := newDB(cfg.DBURL)
	if err != nil {
		return err
	}
	store := sqlLiteStore.NewStore(db)
	inf := infra.NewInfra(db)
	err = inf.LoadDBScheduleItem()
	if err != nil {
		slog.Warn(err.Error())
	}
	srv := server.NewServer(store)

	tgbot.StartBot(srv, inf)

	defer db.Close()
	return nil
}

func newDB(dbURL string) (*sqlx.DB, error) {
	slog.Info("Подключение к DB")
	db, err := sqlx.Open("sqlite", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}
