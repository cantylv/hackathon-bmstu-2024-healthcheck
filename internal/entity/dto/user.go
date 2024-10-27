package dto

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
)

var (
	ErrInvalidUsernameText = errors.New("Никнейм может состоять только из цифр, знака нижнего подчеркивания, символов латинского языка или кириллицы")
	ErrUsernameTooLong     = errors.New("Длина никнейма должна быть не больше 30 символов")
	ErrUsernameTooShort    = errors.New("Длина никнейма должна быть не меньше 2 символов")

	ErrInvalidFirstName  = errors.New("Имя может состоять только из символов латинского языка или кириллицы")
	ErrFirstNameTooLong  = errors.New("Длина имени должна быть не больше 30 символов")
	ErrFirstNameTooShort = errors.New("Длина имени должна быть не меньше 2 символов")

	ErrInvalidWeight = errors.New("Масса тела должна быть положительной")
	ErrInvalidHeight = errors.New("Значение роста должно быть положительным")
	ErrInvalidAge    = errors.New("Значение возраста должно быть положительным")

	ErrInvalidSex      = errors.New("Указан несуществующий пол человека")
	ErrInvalidActivity = errors.New("Такого вида активности не существует")

	ErrInvalidPasswordText = errors.New("Пароль должен содержать как минимум одну цифру и одну заглавную букву")
	ErrPasswordTooLong     = errors.New("Длина пароля должна быть не больше 30 символов")
	ErrPasswordTooShort    = errors.New("Длина пароля должна быть не меньше 8 символов")
)

var (
	ErrInvalidEmail = errors.New("Неверная почта, верный формат username@subdomain.tld, например: student@bmstu.ru")
	ErrInvalidLogin = errors.New("Введен неправильный логин")
)

var (
	usernameRegexp  = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9_]{2,30}$`)
	firstNameRegexp = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]{2,30}$`)
	passwordRegexp  = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
)

// INPUT DATAFLOW
type CreateData struct {
	Username         string  `json:"username"`
	FirstName        string  `json:"first_name"`
	Weight           float32 `json:"weight"`
	Height           int     `json:"height"`
	Age              int     `json:"age"`
	Sex              string  `json:"sex"`
	PhysicalActivity string  `json:"physical_activity"`
	Password         string  `json:"password"`
}

func (h *CreateData) Validate() error {
	// first_name
	err := ValidateUsername(h.Username)
	if err != nil {
		return err
	}

	// weight
	if h.Weight <= 0 {
		return ErrInvalidWeight
	}

	// height
	if h.Height <= 0 {
		return ErrInvalidHeight
	}

	// age
	if h.Age <= 0 {
		return ErrInvalidAge
	}

	// physical_activity
	if _, ok := myconstants.AllowedActivities[h.PhysicalActivity]; !ok {
		return ErrInvalidActivity
	}

	// sex
	if _, ok := myconstants.AllowedHumanSex[h.Sex]; !ok {
		return ErrInvalidSex
	}
	return isPasswordValid(h.Password)
}

func ValidateUsername(username string) error {
	// username
	if len(username) < 2 {
		return ErrUsernameTooShort
	}
	if len(username) > 30 {
		return ErrUsernameTooLong
	}
	if isMatch := usernameRegexp.MatchString(username); !isMatch {
		return ErrInvalidUsernameText
	}
	return nil
}

func isPasswordValid(pwd string) error {
	pwdLen := utf8.RuneCountInString(pwd)
	if pwdLen > 30 {
		return ErrPasswordTooLong
	}
	if pwdLen < 8 {
		return ErrPasswordTooShort
	}

	// проверяем наличие хотя бы одной буквы верхнего регистра
	letterRegex := regexp.MustCompile(`[A-Z]`)
	if !letterRegex.MatchString(pwd) {
		return ErrInvalidPasswordText
	}

	// проверяем наличие хотя бы одной цифры
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(pwd) {
		return ErrInvalidPasswordText
	}

	if !passwordRegexp.MatchString(pwd) {
		return ErrInvalidPasswordText
	}
	return nil
}

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthData) Validate() error {
	countErr := 0
	if isMatch := usernameRegexp.MatchString(h.Login); !isMatch {
		countErr++
	}
	if isMatch := govalidator.IsEmail(h.Login); !isMatch {
		countErr++
	}
	if countErr == 2 {
		return ErrInvalidLogin
	}

	return isPasswordValid(h.Password)
}

// OUTPUT DATAFLOW
type UserWithoutPassword struct {
	ID               string  `json:"id"`
	Email            string  `json:"email"`
	Username         string  `json:"username"`
	FirstName        string  `json:"first_name"`
	Weight           float32 `json:"weight"`
	Height           int     `json:"height"`
	Age              int     `json:"age"`
	Sex              string  `json:"sex"`
	DayCalories      float32 `json:"day_calories"`
	PhysicalActivity string  `json:"physical_activity"`
}
