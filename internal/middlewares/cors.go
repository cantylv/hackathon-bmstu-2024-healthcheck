// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"net/http"
)

// CORS (Cross-Origin Resource Sharing). Настраивает политику доступа различных веб-услуг к нашему серверу.
func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Нужен для Postman | в реальной жизни для версии продукта мы должны устанавливать доменные имена вместо "*".
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
		// Preflight-request обработка.
		if r.Method == http.MethodOptions {
			return
		}
		h.ServeHTTP(w, r)
	})
}
