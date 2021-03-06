package gateway

import (
	"github.com/rs/cors"
	"net/http"
)

func newCorsHandler(srv http.Handler, allowedOrigins []string) http.Handler {
	if len(allowedOrigins) == 0 {
		return srv
	}
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{http.MethodPost, http.MethodGet},
		MaxAge:         600,
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(srv)
}
