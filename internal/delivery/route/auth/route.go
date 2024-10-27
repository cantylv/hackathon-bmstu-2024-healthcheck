package auth

import (
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/delivery/auth"
	rUser "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/repo/user"
	ucAuth "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/usecase/auth"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов для работы с пользователями (получение, удаление, создание).
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(postgresClient)
	usecaseAuth := ucAuth.NewUsecaseLayer(repoUser)
	authHandlerManager := auth.NewAuthHandlerManager(usecaseAuth, logger)
	// ручки, отвечающие за сессию пользователя
	r.HandleFunc("/signup", authHandlerManager.SignUp).Methods("POST")   // регистрация
	r.HandleFunc("/signin", authHandlerManager.SignIn).Methods("POST")   // авторизация
	r.HandleFunc("/signout", authHandlerManager.SignOut).Methods("POST") // деавторизация
}
