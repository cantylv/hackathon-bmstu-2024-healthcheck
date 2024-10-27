package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Init инициализирует цепочку middlewares.
func Init(r *mux.Router, logger *zap.Logger) (h http.Handler) {
	h = JwtVerification(h, logger)
	h = Cors(r)
	h = Recover(h, logger)
	h = Access(h, logger)
	return h
}
