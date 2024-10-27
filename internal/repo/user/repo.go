package user

import (
	"context"
	"fmt"

	ent "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity"
	repoErr "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	GetByUsername(ctx context.Context, username string) (*ent.User, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	DeleteByUsername(ctx context.Context, username string) error
	Create(ctx context.Context, initData *ent.User) (*ent.User, error)
}

var _ Repo = (*RepoLayer)(nil)

type RepoLayer struct {
	dbConn *pgx.Conn
}

// NewRepoLayer возвращает структуру уровня repository. Позволяет работать с пользователем (crd).
func NewRepoLayer(dbConn *pgx.Conn) *RepoLayer {
	return &RepoLayer{
		dbConn: dbConn,
	}
}

var (
	user_fields = "id, username, first_name, weight, height, age, sex, physical_activity, day_calories, password"
)

var (
	sqlRowGetByUsername = fmt.Sprintf(
		`SELECT %s FROM "user" WHERE username=$1`,
		user_fields,
	)
	sqlRowCreateUser = fmt.Sprintf(`
		INSERT INTO "user" (
			username,
			first_name,  
			weight,
			height,
			age,
			sex, 
			physical_activity,
			day_calories,
			password   
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING %s`, user_fields)
)

// GetByUsername позволяет получить пользователя с помощью никнейма.
func (r *RepoLayer) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowGetByUsername, username)
	var u ent.User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.Weight, &u.Height, &u.Age, &u.Sex, &u.PhysicalActivity, &u.DayCalories, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail позволяет получить пользователя с помощью почты пользователя.
func (r *RepoLayer) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowGetByUsername, email)
	var u ent.User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.Weight, &u.Height, &u.Age, &u.Sex, &u.PhysicalActivity, &u.DayCalories, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// DeleteByUsername позволяет удалить пользователя из системы.
func (r *RepoLayer) DeleteByUsername(ctx context.Context, username string) error {
	row, err := r.dbConn.Exec(ctx, `DELETE FROM "user" WHERE username = $1`, username)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return repoErr.ErrNoRowsAffected
	}
	return nil
}

// Create позволяет создать пользователя.
func (r *RepoLayer) Create(ctx context.Context, initData *ent.User) (*ent.User, error) {
	row := r.dbConn.QueryRow(ctx, sqlRowCreateUser,
		initData.Username,
		initData.FirstName,
		initData.Weight,
		initData.Height,
		initData.Age,
		initData.Sex,
		initData.PhysicalActivity,
		initData.DayCalories,
		initData.Password,
	)
	var u ent.User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.Weight, &u.Height, &u.Age, &u.Sex, &u.PhysicalActivity, &u.DayCalories, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
