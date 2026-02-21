package tgbot

import (
	"fmt"
	"log/slog"
	"queue/internal/entity"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (ctl *Controller) registerQueueHandlers(b *tele.Bot) {
	ctl.registerJoinHandler(b)
	ctl.registerLeaveHandler(b)
}

func (ctl *Controller) registerLeaveHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "leave"}, func(c tele.Context) error {
		_ = c.Respond()
		id := c.Sender().ID
		scheduleItemID, _ := strconv.Atoi(c.Data())
		err := ctl.srv.Leave(id, scheduleItemID)
		if err == entity.ErrAlreadyRegistered {
			return c.Respond(&tele.CallbackResponse{Text: "Вы уже покинули очередь"})
		}
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
		}
		var text string
		item, err := ctl.srv.GetItemByID(scheduleItemID)
		if err != nil {
			c.Respond(&tele.CallbackResponse{Text: "не удалось получить пару"})
		}
		if item != nil {
			desc := strings.ReplaceAll(item[0].Description, `\n`, "\n")
			text += fmt.Sprintf("%s\n%s\n", item[0].Name, desc)
		}
		text += fmt.Sprintf("Очередь:\n")
		queue, err := ctl.srv.GetUserByItemID(scheduleItemID)
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "ошибка при получении очереди"})
		}
		i := 1
		for _, it := range queue {
			text += fmt.Sprintf("%d.  %s\n", i, it.Name)
			i++
		}
		return c.Edit(text, ctl.kb.LessonActions(int64(scheduleItemID)))
	})
}

func (ctl *Controller) registerJoinHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "join"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка записи на занятие (join)")
		userID := c.Sender().ID
		scheduleItemID, err := strconv.Atoi(c.Data())
		if err != nil {
			slog.Warn(err.Error())
			return c.Respond(&tele.CallbackResponse{Text: "Некорректный ID занятия"})
		}
		err = ctl.srv.AddUserToItem(userID, scheduleItemID)
		if err == entity.ErrAlreadyRegistered {
			return c.Respond(&tele.CallbackResponse{Text: "пользователь уже записан"})
		}
		if err == entity.ErrUserNotFound {
			return c.Respond(&tele.CallbackResponse{Text: "Сначала нажми /start и введи имя"})
		}
		if err == entity.ErrScheduleNotFound {
			return c.Respond(&tele.CallbackResponse{Text: "Пара устарела, открой расписание заново"})
		}
		if err != nil {
			slog.Warn(err.Error())
			return c.Respond(&tele.CallbackResponse{Text: "Не удалось записать ❌ попробуйте ввести имя"})
		}

		var text string
		item, err := ctl.srv.GetItemByID(scheduleItemID)
		if err != nil {
			c.Respond(&tele.CallbackResponse{Text: "не удалось получить пару"})
		}
		if item != nil {
			desc := strings.ReplaceAll(item[0].Description, `\n`, "\n")
			text += fmt.Sprintf("%s\n%s\n", item[0].Name, desc)
		}
		text += fmt.Sprintf("Очередь:\n")
		queue, err := ctl.srv.GetUserByItemID(scheduleItemID)
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "ошибка при получении очереди"})
		}
		i := 1
		for _, it := range queue {
			text += fmt.Sprintf("%d.  %s\n", i, it.Name)
			i++
		}
		c.Respond(&tele.CallbackResponse{Text: "Записал ✅"})
		return c.Edit(text, ctl.kb.LessonActions(int64(scheduleItemID)))
	})
}
