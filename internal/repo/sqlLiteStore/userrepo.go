package sqlLiteStore

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

// NewUser создание пользователя
func (u UserRepo) NewUser(id int64, name string) error {
	slog.Debug("запрос в DB создание пользователя NewUser")
	_, err := u.db.Exec(`insert into users (id,name)
					 		   values ( ?,?)`, id, name)
	return err
}

// RenameUser обновление имени пользователя
func (u UserRepo) RenameUser(id int64, newName string) error {
	slog.Debug("запрос в DB смена имени пользователю RenameUser")
	_, err := u.db.Exec(`update users set name = ? where id = ?`, newName, id)
	return err
}

// Exists 1 есть пользователь 0 нет
func (u UserRepo) Exists(id int64) (bool, error) {
	slog.Debug("запрос DB проверка существует ли пользователь Exists")
	var exists bool
	err := u.db.Get(&exists, "select exists(select 1 from users where id = ?)", id)
	if err != nil {
		return false, err
	}
	return exists, nil
}
