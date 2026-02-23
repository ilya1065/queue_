package tgbot

import (
	"log"
	"log/slog"
	"os"
	"queue/internal/infra"
	"queue/internal/server"
	"time"

	"gopkg.in/telebot.v4"
)

// StartBot запуск Telegram-бота и регистрация всех роутов
func StartBot(srv *server.Server, inf *infra.Infra) {
	slog.Info("Старт бота")
	if os.Getenv("TG_KEY") == "" {
		log.Fatal("TG_KEY is not set")
	}

	pref := telebot.Settings{
		Token:  os.Getenv("TG_KEY"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	ctl := NewController(srv, inf)
	ctl.RegisterRoutes(b)
	slog.Info("Бот готов, запускаю polling")
	b.Start()
}
