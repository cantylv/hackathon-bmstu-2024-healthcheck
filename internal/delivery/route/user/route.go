package user

import (
	dUser "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/delivery/user"
	rUser "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/repo/user"
	ucUser "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHandlers инициализирует обработчики запросов для работы с пользователями (получение, удаление, создание).
func InitHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(postgresClient)
	ucUser := ucUser.NewUsecaseLayer(repoUser)
	userHandlerManager := dUser.NewUserHandlerManager(ucUser, logger)
	// ручки, отвечающие за получение и удаление пользователя
	r.HandleFunc("/users", userHandlerManager.Read).Methods("GET")        // чтение данных пользователя
	r.HandleFunc("/users", userHandlerManager.Delete).Methods("DELETE")   // удаление пользователя
	r.HandleFunc("/users/weight", userHandlerManager.UpdateWeight).Methods("PUT") // обновление массы тела
}
