package tgbot

import (
	"fmt"
	"log/slog"
	"queue/internal/entity"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Keyboards struct{}

func NewKeyboards() *Keyboards { return &Keyboards{} }

func (k *Keyboards) WeekMenu() *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	btnCur := m.Data("Текущая неделя", "week_current", "")
	btnNext := m.Data("Следующая неделя", "week_next", "")
	backMain := m.Data("Назад", "back", "")
	m.Inline(m.Row(btnCur, btnNext), m.Row(backMain))
	return m
}

func (k *Keyboards) MainMenu(id int64) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	slog.Debug(fmt.Sprintf("MainMenu %d", id))
	name := m.Data("Вести имя", "name", "")
	record := m.Data("Записаться", "record", "")
	if id == 8141813763 {
		adminMenu := m.Data("админ панель", "admin_menu", "")
		m.Inline(m.Row(name), m.Row(record), m.Row(adminMenu))
	} else {
		m.Inline(m.Row(name), m.Row(record))
	}
	return m
}

func (k *Keyboards) AdminMenu() *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	reload := m.Data("Обновить расписание", "reload", "")
	back := m.Data("назад", "back", "")
	m.Inline(m.Row(reload), m.Row(back))
	return m
}

func (k *Keyboards) DaysMenu(weekStart time.Time) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	labels := []string{"ПН", "ВТ", "СР", "ЧТ", "ПТ", "СБ", "ВС"}

	var rows []tele.Row
	var row []tele.Btn

	for i := 0; i < 7; i++ {
		day := weekStart.AddDate(0, 0, i)
		btn := m.Data(labels[i], "day", day.Format("2006-01-02"))
		row = append(row, btn)
		if len(row) == 4 {
			rows = append(rows, m.Row(row...))
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, m.Row(row...))
	}

	back := m.Data("⬅️ Назад", "back_week", "")
	rows = append(rows, m.Row(back))

	m.Inline(rows...)
	return m
}

func (k *Keyboards) LessonsMenu(items []entity.ScheduleItem) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	var rows []tele.Row
	var row []tele.Btn
	for i, it := range items {
		btn := m.Data(fmt.Sprintf("%s %d", it.Name, i+1), "lesson", fmt.Sprintf("%d", it.Id))
		row = append(row, btn)

		//if len(row) == 5 {
		rows = append(rows, m.Row(row...))
		row = nil
		//}

	}
	//if len(row) > 0 {
	//	rows = append(rows, m.Row(row...))
	//}
	rows = append(rows, m.Row(m.Data("⬅️ Назад", "back_days", "")))
	m.Inline(rows...)
	return m
}

func (k *Keyboards) LessonActions(scheduleItemID int64) *tele.ReplyMarkup {
	m := &tele.ReplyMarkup{}
	join := m.Data("✅ Записаться", "join", fmt.Sprintf("%d", scheduleItemID))
	leave := m.Data("Покинуть Очередь", "leave", fmt.Sprintf("%d", scheduleItemID))
	back := m.Data("⬅️ Назад", "back_lessons", "")
	m.Inline(
		m.Row(join),
		m.Row(leave),
		m.Row(back),
	)
	return m
}
