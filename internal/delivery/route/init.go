package route

import (
	"net/http"

	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/delivery/route/auth"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/delivery/route/user"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/middlewares"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// InitHTTPHandlers инициализирует обработчики запросов, а также добавляет цепочку middlewares в обработку запроса.
func InitHTTPHandlers(r *mux.Router, postgresClient *pgx.Conn, logger *zap.Logger) http.Handler {
	s := r.PathPrefix("/api/v1").Subrouter()
	user.InitHandlers(s, postgresClient, logger)
	auth.InitHandlers(s, postgresClient, logger)
	return middlewares.Init(s, logger)
}
