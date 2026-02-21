package tgbot

import (
	"log/slog"
	"queue/internal/infra"
	"queue/internal/server"
	"sync"
	"time"

	tele "gopkg.in/telebot.v4"
)

var (
	errorRetrievingSchedule = "Ошибка получения расписания"
)

type Controller struct {
	kb  *Keyboards
	srv *server.Server
	inf *infra.Infra
	// Минимальный контекст в памяти (потом заменишь на repo)
	userWeekStart sync.Map // key int64 -> time.Time
	userDay       sync.Map // key int64 -> time.Time
	waitName      sync.Map
}

func NewController(srv *server.Server, inf *infra.Infra) *Controller {
	return &Controller{
		inf: inf,
		kb:  NewKeyboards(),
		//userWeekStart: make(map[int64]time.Time),
		//userDay:       make(map[int64]time.Time),
		srv: srv,
	}
}

func weekStart(d time.Time) time.Time {
	slog.Debug("получение начала недели")
	wd := int(d.Weekday())
	if wd == 0 {
		wd = 7 // Sunday -> 7
	}
	base := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	return base.AddDate(0, 0, -(wd - 1))
}

// RegisterRoutes — вызывай из main после создания bot
func (ctl *Controller) RegisterRoutes(b *tele.Bot) {
	slog.Info("регистрация роутера")
	ctl.registerUserHandlers(b)
	ctl.registerScheduleHandlers(b)
	ctl.registerQueueHandlers(b)
	ctl.registerAdminHandlers(b)
}
