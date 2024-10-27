package myerrors

import "errors"

var (
	ErrInternal             = errors.New("Внутренняя ошибка сервера, пожалуйста, попробуйте немного позже")
	ErrNoRequestIdInContext = errors.New("Отсутствует request_id в контексте объекта запроса")

	ErrInvalidJwt = errors.New("Невалидный jwt_token")

	ErrIncorrectPwdOrLogin  = errors.New("Неверный пароль или логин")
	ErrUserAlreadyExist     = errors.New("Пользователь с таким никнеймом уже существует")
	ErrUserNotExist         = errors.New("Пользователь с таким никнеймом не существует")
	ErrAlreadyRegistered    = errors.New("Вы уже зарегистрированы")
	ErrAlreadyAuthenticated = errors.New("Вы уже авторизованы")
	ErrNotAuthenticated     = errors.New("Вы не авторизованы")
	ErrInvalidData          = errors.New("Вы ввели неправильные данные")
)

var (
	ErrNoRowsAffected = errors.New("no rows were affected")
)
