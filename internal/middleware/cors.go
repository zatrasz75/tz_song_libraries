package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

// SetHeader устанавливает заголовок Access-Control-Allow-Origin в ответе.
func SetHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		next.ServeHTTP(w, r)
	})
}

// CORS создает новый экземпляр middleware с заданными настройками.
func CORS(allowOrigins []string) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization", "x-token"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"*"},
		MaxAge:           43200, // 12 часов
	})

	return c.Handler
}
