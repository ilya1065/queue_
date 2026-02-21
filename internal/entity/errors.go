package entity

import "errors"

var (
	ErrAlreadyRegistered = errors.New("пользователь уже записан")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrScheduleNotFound  = errors.New("пара не найдена")
)
