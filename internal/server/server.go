package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"queue/internal/config"
	"queue/internal/entity"
	"queue/internal/repo"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

type Server struct {
	store repo.Store
	cfg   *config.Config
}

func NewServer(store repo.Store, cfg *config.Config) *Server {
	slog.Info("создание сервера")
	return &Server{
		store: store,
		cfg:   cfg,
	}
}

// Leave выход из очереди
func (s *Server) Leave(userID int64, scheduleID int) error {
	err := s.store.Record().DeleteRecord(userID, scheduleID)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	return nil
}

// GetItemByTime возвращает срез пар на день даты start
func (s *Server) GetItemByTime(start time.Time) ([]entity.ScheduleItem, error) {
	slog.Debug("сервис GetItemByTime")
	end := start.AddDate(0, 0, 1)
	items, err := s.store.ScheduleItem().GetItemByTime(start, end)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Server) UpdateUsers(name string, id int64) error {
	slog.Debug("сервис UpdateUsers")
	exists, err := s.store.User().Exists(id)
	if err != nil {
		return err
	}
	if exists {
		return s.store.User().RenameUser(id, name)
	}
	return s.store.User().NewUser(id, name)

}

func (s *Server) ExistsUser(id int64) (bool, error) {
	return s.store.User().Exists(id)
}
func (s *Server) AddUserToItem(id int64, scheduleItemID int) error {
	err := s.store.Record().AddUserToItem(id, scheduleItemID)
	return err
}

func (s *Server) GetUserByItemID(id int) ([]entity.User, error) {
	users, err := s.store.Record().GetUserByItemID(id)
	return users, err
}

func (s *Server) GetItemByID(id int) (*entity.ScheduleItem, error) {
	return s.store.ScheduleItem().GetItemByID(id)
}

func (s *Server) UpdateScheduleForNextTwoWeeks() error {
	slog.Debug("Запрос к api обновление расписания")
	type Response struct {
		PageProps struct {
			ScheduleLoadInfo []struct {
				ID          int    `json:"id"`
				Title       string `json:"title"`
				ICalContent string `json:"iCalContent"`
				ICalLink    string `json:"iCalLink"`
			} `json:"scheduleLoadInfo"`
		} `json:"pageProps"`
	}
	start := timeOfTwoWeekNext()
	end := start.AddDate(1, 0, 0)
	resp, err := http.Get(s.cfg.SchedulerURL)
	if err != nil {
		slog.Warn(err.Error())
		return errors.New("ошибка выполнения запроса к API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	var body Response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		slog.Warn(err.Error())
		return errors.New("ошибка при декодирование данных")
	}
	if len(body.PageProps.ScheduleLoadInfo) == 0 {
		return errors.New("ошибка получения расписания")
	}
	r := strings.NewReader(body.PageProps.ScheduleLoadInfo[0].ICalContent)
	c := gocal.NewParser(r)
	c.Start = &start
	c.End = &end
	err = c.Parse()
	if err != nil {
		slog.Warn(err.Error())
		return errors.New("ошибка парсинга расписания")
	}
	var items []entity.ScheduleItem
	for _, item := range c.Events {
		items = append(items, entity.ScheduleItem{
			Name:        item.Summary,
			Description: item.Description,
			StartDate:   item.Start,
			EndDate:     item.End,
			ExternalID:  item.Uid,
		})
	}
	err = s.store.ScheduleItem().UpdateScheduleForTime(items, start)
	if err != nil {
		slog.Warn(err.Error())
		return errors.New(fmt.Sprintf("Ошибка: %v", err))
	}
	return nil
}

func timeOfTwoWeekNext() time.Time {
	t := time.Now()
	year, month, day := t.Date()
	location := t.Location()
	t = time.Date(year, month, day, 1, 0, 0, 0, location)
	weekday := int(t.Weekday())
	offset := (weekday + 6) % 7
	startOfThisWeek := t.AddDate(0, 0, -offset)
	return startOfThisWeek.AddDate(0, 0, 14)
}
