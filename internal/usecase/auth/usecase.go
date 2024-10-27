package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/asaskevich/govalidator"
	ent "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/repo/user"
	f "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/functions"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
)

type Usecase interface {
	SignUp(ctx context.Context, signUpData *dto.CreateData) (*ent.User, error)
	SignIn(ctx context.Context, authData *dto.AuthData) (*ent.User, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoUser user.Repo
}

// NewUsecaseLayer возращает структуру уровня usecase для работы с пользователями.
func NewUsecaseLayer(repoUser user.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoUser: repoUser,
	}
}

// SignUp регистрирует пользователя.
func (u *UsecaseLayer) SignUp(ctx context.Context, authData *dto.CreateData) (*ent.User, error) {
	// проверяем, существует ли пользователь c таким некнеймом
	// если да, то возвращаем ошибку
	uDB, err := u.repoUser.GetByUsername(ctx, authData.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if uDB != nil {
		return nil, me.ErrUserAlreadyExist
	}
	// получаем хэшированный пароль вместе с солью
	hashedPassword, err := f.GetHashedPassword(authData.Password)
	if err != nil {
		return nil, err
	}
	dayCalories := f.GetDayCalories(authData)
	userNew, err := u.repoUser.Create(ctx, newUserFromSignUpForm(authData, hashedPassword, dayCalories))
	if err != nil {
		return nil, err
	}
	return userNew, nil
}

// SignIn авторизует пользователя.
func (u *UsecaseLayer) SignIn(ctx context.Context, authData *dto.AuthData) (*ent.User, error) {
	var dbUser *ent.User
	loginIsEmail := govalidator.IsEmail(authData.Login)
	// если авторизация была по почте, то:
	if loginIsEmail {
		uDB, err := u.repoUser.GetByEmail(ctx, authData.Login)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, me.ErrIncorrectPwdOrLogin
			}
			return nil, err
		}
		dbUser = uDB
	} else {
		uDB, err := u.repoUser.GetByUsername(ctx, authData.Login)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, me.ErrIncorrectPwdOrLogin
			}
			return nil, err
		}
		dbUser = uDB
	}
	if !f.IsPasswordsEqual(authData.Password, dbUser.Password) {
		return nil, me.ErrIncorrectPwdOrLogin
	}
	return dbUser, nil
}
