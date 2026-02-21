package tgbot

import (
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func (ctl *Controller) registerAdminHandlers(b *tele.Bot) {
	ctl.registerAdminMenuHandler(b)
	ctl.registerReloadHandler(b)
}

func (ctl *Controller) registerAdminMenuHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "admin_menu"}, func(c tele.Context) error {
		_ = c.Respond()
		return c.Edit("Меню", ctl.kb.AdminMenu())
	})
}

func (ctl *Controller) registerReloadHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "reload"}, func(c tele.Context) error {
		_ = c.Respond()
		err := ctl.inf.LoadDBScheduleItem()
		if err != nil {
			slog.Warn(err.Error())
			return c.Edit(fmt.Sprintf("не полусилось(((\n Меню:"), ctl.kb.AdminMenu())
		}
		return c.Respond(&tele.CallbackResponse{Text: "ОК"})
	})
}
