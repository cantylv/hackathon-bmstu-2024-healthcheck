package main

import (
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/config"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/app"
	"go.uber.org/zap"
)

// main точка старта приложения
func main() {
	logger := zap.Must(zap.NewProduction())
	config.Read("./config/config.yaml", logger)
	app.Run(logger)
}
