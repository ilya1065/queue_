package sqlLiteStore

import (
	"log/slog"
	"queue/internal/repo"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db               *sqlx.DB
	userRepo         *UserRepo
	scheduleItemRepo *ScheduleItemRepo
	recordRepo       *RecordRepo
}

func NewStore(db *sqlx.DB) *Store {
	slog.Info("создание хранилища")
	return &Store{
		db: db,
	}
}

func (s *Store) User() repo.UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}
	s.userRepo = &UserRepo{
		db: s.db,
	}
	return s.userRepo

}

func (s *Store) ScheduleItem() repo.ScheduleRepo {
	if s.scheduleItemRepo != nil {
		return s.scheduleItemRepo
	}
	s.scheduleItemRepo = &ScheduleItemRepo{
		db: s.db,
	}
	return s.scheduleItemRepo
}

func (s *Store) Record() repo.RecordsRepo {
	if s.recordRepo != nil {
		return s.recordRepo
	}
	s.recordRepo = &RecordRepo{
		db: s.db,
	}
	return s.recordRepo
}
