package web

import (
	"fmt"
	"net/http"
)

func (app *App) recoverPanic(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection:", "close")
				var errorMessage string
				if e, ok := err.(error); ok {
					errorMessage = e.Error()
					app.logger.Error(errorMessage)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				} else {
					errorMessage = fmt.Sprintf("panic: %v", err)
					app.logger.Error(errorMessage)
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
}
func (app *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", ("*"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
