package auth

import (
	ent "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
)

func getUserWithoutPassword(user *ent.User) *dto.UserWithoutPassword {
	return &dto.UserWithoutPassword{
		ID:               user.ID,
		Email:            user.Email,
		FirstName:        user.FirstName,
		Username:         user.Username,
		Weight:           user.Weight,
		Height:           user.Height,
		Age:              user.Age,
		Sex:              user.Sex,
		DayCalories:      user.DayCalories,
		PhysicalActivity: user.PhysicalActivity,
	}
}
