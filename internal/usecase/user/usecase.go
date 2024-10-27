package user

import (
	"context"
	"database/sql"
	"errors"

	ent "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/repo/user"
	f "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/functions"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
)

type Usecase interface {
	Read(ctx context.Context, username string) (*ent.User, error)
	Delete(ctx context.Context, username string) error
	UpdateWeight(ctx context.Context, weight float32, username string) (*ent.User, error)
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

// Read возвращает данные о пользователе.
func (u *UsecaseLayer) Read(ctx context.Context, username string) (*ent.User, error) {
	uDB, err := u.repoUser.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	return uDB, nil
}

// Delete удаляет пользователя из системы.
func (u *UsecaseLayer) Delete(ctx context.Context, username string) error {
	// проверка существования пользователя
	_, err := u.repoUser.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return me.ErrUserNotExist
		}
		return err
	}
	return u.repoUser.DeleteByUsername(ctx, username)
}

func (u *UsecaseLayer) UpdateWeight(ctx context.Context, weight float32, username string) (*ent.User, error) {
	// проверка существования пользователя
	uDB, err := u.repoUser.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, me.ErrUserNotExist
		}
		return nil, err
	}
	// посчитаем новое значение КК
	newDayCalories := f.GetDayCalories(newCreateDataFromUser(uDB, weight))
	return u.repoUser.UpdateWeight(ctx, weight, newDayCalories, username)
}
