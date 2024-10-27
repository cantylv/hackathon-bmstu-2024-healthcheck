package functions

import (
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
)

func GetDayCalories(usr *dto.CreateData) float64 {
	var bmr float64
	// Приводим к float64 для всех вычислений
	weight := float64(usr.Weight)
	height := float64(usr.Height)
	age := float64(usr.Age)

	if usr.Sex == "F" {
		bmr = (10 * weight) + (6.25 * height) - (5 * age) - 161
	} else {
		bmr = (10 * weight) + (6.25 * height) - (5 * age) + 5
	}

	// Умножаем BMR на коэффициент активности
	caloriesNeeded := bmr * float64(mc.AllowedActivities[usr.PhysicalActivity])
	return caloriesNeeded
}
