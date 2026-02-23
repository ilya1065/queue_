package tgbot

import (
	"log/slog"
	"strings"

	tele "gopkg.in/telebot.v4"
)

// registerUserHandlers регистрация пользовательских роутов
func (ctl *Controller) registerUserHandlers(b *tele.Bot) {
	ctl.registerStartHandler(b)
	ctl.registerNameHandler(b)
	ctl.registerTextHandler(b)
	ctl.registerBackHandler(b)
}

// registerNameHandler обработчик кнопки смены имени
func (ctl *Controller) registerNameHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "name"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка name")
		ctl.waitName.Store(c.Sender().ID, true)
		return c.Edit("Введите новое имя:")
	})
}

// registerStartHandler обработчик команды /start
func (ctl *Controller) registerStartHandler(b *tele.Bot) {
	b.Handle("/start", func(c tele.Context) error {
		slog.Debug("кнопка start")
		id := c.Sender().ID
		exist, err := ctl.srv.ExistsUser(id)
		if err != nil {
			return err
		}
		if exist {
			return c.Send("Меню", ctl.kb.MainMenu(id))
		}
		ctl.waitName.Store(id, true)
		return c.Send("Введите ваше имя (одно сообщение):")
	})
}

// registerTextHandler обработчик текстового ввода имени
func (ctl *Controller) registerTextHandler(b *tele.Bot) {
	b.Handle(tele.OnText, func(c tele.Context) error {
		slog.Debug("принимаем текст Имени")
		id := c.Sender().ID
		if _, waiting := ctl.waitName.Load(id); !waiting {
			return nil
		}
		name := strings.TrimSpace(c.Text())
		if name == "" {
			return c.Send("Имя не может быть пустым. Введите ещё раз:")
		}
		if len([]rune(name)) > 40 {
			return c.Send("Слишком длинное имя. Введите покороче:")
		}

		if err := ctl.srv.UpdateUsers(name, id); err != nil {
			return err
		}

		ctl.waitName.Delete(id)
		return c.Send("Отлично! Меню:", ctl.kb.MainMenu(id))
	})
}

// registerBackHandler обработчик возврата в главное меню
func (ctl *Controller) registerBackHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "back"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка меню")
		id := c.Sender().ID
		return c.Edit("Меню", ctl.kb.MainMenu(id))
	})
}
