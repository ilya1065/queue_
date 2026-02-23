package tgbot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

// registerScheduleHandlers регистрация роутов расписания.
func (ctl *Controller) registerScheduleHandlers(b *tele.Bot) {
	ctl.registerRecordHandler(b)
	ctl.registerWeekCurrentHandler(b)
	ctl.registerWeekNextHandler(b)
	ctl.registerBackWeekHandler(b)
	ctl.registerDayHandler(b)
	ctl.registerBackDaysHandler(b)
	ctl.registerLessonHandler(b)
	ctl.registerBackLessonsHandler(b)
}

// registerRecordHandler вход в сценарий записи на пару
func (ctl *Controller) registerRecordHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "record"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка записи (record)")
		return c.Edit("Выбери неделю:", ctl.kb.WeekMenu())
	})
}

// registerWeekCurrentHandler выбор текущей недели
func (ctl *Controller) registerWeekCurrentHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "week_current"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка текущей недели")
		ws := weekStart(time.Now())
		return c.Edit("Выбери день:", ctl.kb.DaysMenu(ws))
	})
}

// registerWeekNextHandler выбор следующей недели
func (ctl *Controller) registerWeekNextHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "week_next"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка следующей недели")
		ws := weekStart(time.Now()).AddDate(0, 0, 7)
		return c.Edit("Выбери день:", ctl.kb.DaysMenu(ws))
	})
}

// registerDayHandler выбор дня и загрузка пар на выбранную дату
func (ctl *Controller) registerDayHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "day"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка выбора дня")
		cb := c.Callback()
		slog.Debug("CB",
			"unique", cb.Unique,
			"data", c.Data(),
			"msg_id", cb.Message.ID,
			"from", c.Sender().ID,
			"cb_id", cb.ID,
		)

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Ошибка часового пояса"})
		}
		day, err := time.ParseInLocation("2006-01-02", c.Data(), loc)
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Некорректная дата"})
		}

		ctl.userDay.Store(c.Sender().ID, day)
		ctl.userWeekStart.Store(c.Sender().ID, weekStart(day))

		items, err := ctl.srv.GetItemByTime(day)
		for _, it := range items {
			slog.Debug("item", "id", it.Id, "name", it.Name, "time", it.StartDate)
		}
		if err != nil {
			fmt.Println(err)
			return c.Respond(&tele.CallbackResponse{Text: errorRetrievingSchedule})
		}

		text := fmt.Sprintf("Расписание на %s: \n", day.Format("2006-01-02"))
		if len(items) == 0 {
			text += "Нет занятий.\n"
		} else {
			text += "\n\nВыбери пару:"
		}
		return c.Edit(text, ctl.kb.LessonsMenu(items))
	})
}

// registerLessonHandler выбор пары и показ текущей очереди
func (ctl *Controller) registerLessonHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "lesson"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("выбор пары lesson")
		scheduleItemID, err := strconv.Atoi(c.Data())
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Некорректный ID пары"})
		}
		var text string
		item, err := ctl.srv.GetItemByID(scheduleItemID)
		if err != nil {
			c.Respond(&tele.CallbackResponse{Text: "не удалось получить пару"})
		}
		if item != nil {
			desc := strings.ReplaceAll(item.Description, `\n`, "\n")
			text += fmt.Sprintf("%s\n%s\n", item.Name, desc)
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

// registerBackWeekHandler возврат к выбору недели
func (ctl *Controller) registerBackWeekHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "back_week"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка назад к выбору недели")
		return c.Edit("Выбери неделю:", ctl.kb.WeekMenu())
	})
}

// registerBackDaysHandler возврат к выбору дня
func (ctl *Controller) registerBackDaysHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "back_days"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка назад к выбору дня")
		v, ok := ctl.userWeekStart.Load(c.Sender().ID)
		var ws time.Time
		if !ok {
			ws = weekStart(time.Now())
		} else {
			ws = v.(time.Time)
		}
		return c.Edit("Выбери день:", ctl.kb.DaysMenu(ws))
	})
}

// registerBackLessonsHandler возврат к списку пар выбранного дня
func (ctl *Controller) registerBackLessonsHandler(b *tele.Bot) {
	b.Handle(&tele.InlineButton{Unique: "back_lessons"}, func(c tele.Context) error {
		_ = c.Respond()
		slog.Debug("кнопка назад к парам")

		v, ok := ctl.userDay.Load(c.Sender().ID)
		if !ok {
			return c.Edit("Выбери неделю:", ctl.kb.WeekMenu())
		}

		day, ok := v.(time.Time)
		if !ok {
			slog.Error("userDay contains non-time value")
			return c.Edit("Выбери неделю:", ctl.kb.WeekMenu())
		}

		items, err := ctl.srv.GetItemByTime(day)
		for _, it := range items {
			slog.Debug("item", "id", it.Id, "name", it.Name)
		}
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: errorRetrievingSchedule})
		}

		return c.Edit(
			fmt.Sprintf("Расписание на %s\n", day.Format("02.01.2006")),
			ctl.kb.LessonsMenu(items),
		)
	})
}
