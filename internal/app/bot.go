package app

import (
	"log/slog"
	"queue/internal/config"
	"queue/internal/infra"
	"queue/internal/repo/sqlLiteStore"
	"queue/internal/server"
	"queue/internal/tgbot"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Start сборка зависимостей и запуск приложения
func Start(cfg *config.Config) error {
	slog.Info("Старт сервиса")
	db, err := newDB(cfg.DBURL)
	if err != nil {
		return err
	}
	store := sqlLiteStore.NewStore(db, cfg)
	inf := infra.NewInfra(db, cfg)
	start := time.Now().Add(-1 * 90 * 24 * time.Hour)
	end := time.Now().Add(1 * 365 * 24 * time.Hour)
	err = inf.LoadDBScheduleItem(start, end)
	if err != nil {
		slog.Warn(err.Error())
	}
	if err != nil {
		slog.Warn(err.Error())
	}
	srv := server.NewServer(store, cfg)

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
