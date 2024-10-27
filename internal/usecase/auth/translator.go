package auth

import (
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
)

func newUserFromSignUpForm(data *dto.CreateData, hashedPassword string, dayCalories float64) *entity.User {
	return &entity.User{
		Username:         data.Username,
		FirstName:        data.FirstName,
		Weight:           data.Weight,
		Height:           data.Height,
		Age:              data.Age,
		Sex:              data.Sex,
		PhysicalActivity: data.PhysicalActivity,
		DayCalories:      float32(dayCalories),
		Password:         hashedPassword,
	}
}
