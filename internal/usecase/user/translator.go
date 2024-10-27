package user

import (
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
)

func newCreateDataFromUser(data *entity.User, weight float32) *dto.CreateData {
	return &dto.CreateData{
		Weight:           weight,
		Height:           data.Height,
		Age:              data.Age,
		Sex:              data.Sex,
		PhysicalActivity: data.PhysicalActivity,
	}
}
