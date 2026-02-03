package tgbot

import (
	"log"
	"log/slog"
	"queue/internal/config"
	"queue/internal/infra"
	"queue/internal/server"
	"time"

	"gopkg.in/telebot.v4"
)

func StartBot(cfg *config.Config, srv *server.Server, inf *infra.Infra) {
	slog.Info("Старт бота")
	pref := telebot.Settings{
		Token:  cfg.TGKey,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	ctl := NewController(srv, inf)
	ctl.RegisterRoutes(b)
	b.Handle("/hello", func(c telebot.Context) error {
		return c.Send("Hello")
	})
	b.Start()
}
