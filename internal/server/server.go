package server

import (
	"log/slog"
	"queue/internal/entity"
	"queue/internal/repo"
	"time"
)

type Server struct {
	store repo.Store
}

func NewServer(store repo.Store) *Server {
	slog.Info("создание сервера")
	return &Server{
		store: store,
	}
}

func (s *Server) Leave(userID int64, scheduleID int) error {
	err := s.store.Record().DeleteRecord(userID, scheduleID)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	return nil
}

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
	} else {
		return s.store.User().NewUser(id, name)
	}
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

func (s *Server) GetItemByID(id int) ([]entity.ScheduleItem, error) {
	return s.store.ScheduleItem().GetItemByID(id)
}
