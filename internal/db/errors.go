package db

import "errors"

var (
	userNotFoundError             = errors.New("Не удалось найти данного пользователя")
	incorrectLoginOrPasswordError = errors.New("Неверный логин или пароль")
	unknownError                  = errors.New("Неизвестная ошибка, попробуйте позже")
)
