package tgbot

import (
	"fmt"
	"log/slog"
	"queue/internal/entity"
	"queue/internal/infra"
	"queue/internal/server"
	"strconv"
	"strings"
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
	// /start
	//b.Handle("/start", func(c tele.Context) error {
	//	return c.Send("Выбери неделю:", ctl.kb.WeekMenu())
	//})
	b.Handle(&tele.InlineButton{Unique: "name"}, func(c tele.Context) error {
		slog.Debug("кнопка name")
		ctl.waitName.Store(c.Sender().ID, true)
		return c.Edit("Введите новое имя:")
	})

	//b.Handle(&tele.InlineButton{Unique: "name"}, func(c tele.Context) error {
	//
	//})

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

		// сохраняем
		if err := ctl.srv.UpdateUsers(name, id); err != nil {
			return err
		}

		ctl.waitName.Delete(id)
		return c.Send("Отлично! Меню:", ctl.kb.MainMenu(id))
	})

	b.Handle(&tele.InlineButton{Unique: "back"}, func(c tele.Context) error {
		slog.Debug("кнопка меню")
		id := c.Sender().ID
		return c.Edit("Меню", ctl.kb.MainMenu(id))
	})
	b.Handle(&tele.InlineButton{Unique: "record"}, func(c tele.Context) error {
		slog.Debug("кнопка записи (record)")
		return c.Edit("Выбери неделю:", ctl.kb.WeekMenu())
	})
	// выбор недели
	b.Handle(&tele.InlineButton{Unique: "week_current"}, func(c tele.Context) error {
		slog.Debug("кнопка текущей недели")
		ws := weekStart(time.Now())

		//	items, err := ctl.srv.GetItemByTime(time.Now(), time.Now().AddDate(1, 0, 0))
		//	fmt.Println(items, err)
		return c.Edit("Выбери день:", ctl.kb.DaysMenu(ws))
	})

	b.Handle(&tele.InlineButton{Unique: "week_next"}, func(c tele.Context) error {
		slog.Debug("кнопка следующей недели")
		ws := weekStart(time.Now()).AddDate(0, 0, 7)

		return c.Edit("Выбери день:", ctl.kb.DaysMenu(ws))
	})

	// выбор дня
	b.Handle(&tele.InlineButton{Unique: "day"}, func(c tele.Context) error {
		slog.Debug("кнопка выбора дня")
		day, err := time.Parse("2006-01-02", c.Data())
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: "Некорректная дата"})
		}

		ctl.userDay.Store(c.Sender().ID, day)                  //[c.Sender().ID] = day
		ctl.userWeekStart.Store(c.Sender().ID, weekStart(day)) //[c.Sender().ID] = weekStart(day)

		items, err := ctl.srv.GetItemByTime(day)
		fmt.Println(items, day, err)
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

	// выбор пары
	b.Handle(&tele.InlineButton{Unique: "lesson"}, func(c tele.Context) error {
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
			text += fmt.Sprintf("%d.  %s\n\n", i, it.Name)
			i++
		}
		return c.Edit(text, ctl.kb.LessonActions(int64(scheduleItemID)))
	})

	// join
	b.Handle(&tele.InlineButton{Unique: "join"}, func(c tele.Context) error {
		// scheduleItemID из Data
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
		if err != nil {
			slog.Warn(err.Error())
			return c.Respond(&tele.CallbackResponse{Text: "Не удалось записать ❌ попробуйте ввести имя"})
		}
		//ctl.srv.AddUserToItem()
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

	// back: к выбору недели
	b.Handle(&tele.InlineButton{Unique: "back_week"}, func(c tele.Context) error {
		slog.Debug("кнопка назад к выбору недели")
		return c.Edit("Выбери неделю:", ctl.kb.WeekMenu())
	})

	// back: к дням
	b.Handle(&tele.InlineButton{Unique: "back_days"}, func(c tele.Context) error {
		slog.Debug("кнопка назад к выбору дня")
		v, ok := ctl.userWeekStart.Load(c.Sender().ID) //[c.Sender().ID]
		var ws time.Time
		if !ok {
			ws = weekStart(time.Now())
		} else {
			ws = v.(time.Time)
		}
		return c.Edit("Выбери день:", ctl.kb.DaysMenu(ws))
	})

	b.Handle(&tele.InlineButton{Unique: "back_lessons"}, func(c tele.Context) error {
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
		if err != nil {
			return c.Respond(&tele.CallbackResponse{Text: errorRetrievingSchedule})
		}

		return c.Edit(
			fmt.Sprintf("Расписание на %s\n", day.Format("02.01.2006")),
			ctl.kb.LessonsMenu(items),
		)
	})

	b.Handle(&tele.InlineButton{Unique: "admin_menu"}, func(c tele.Context) error {
		return c.Edit("Меню", ctl.kb.AdminMenu())
	})

	b.Handle(&tele.InlineButton{Unique: "reload"}, func(c tele.Context) error {
		err := ctl.inf.LoadDBScheduleItem()
		if err != nil {
			slog.Warn(err.Error())
			return c.Edit(fmt.Sprintf("не полусилось(((\n Меню:"), ctl.kb.AdminMenu())
		}
		return c.Respond(&tele.CallbackResponse{Text: "ОК"})
	})
}
