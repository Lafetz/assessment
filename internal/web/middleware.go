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

// func (app *App) cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		w.Header().Set("Vary", "Origin")

// 		w.Header().Set("Vary", "Access-Control-Request-Method")
// 		origin := r.Header.Get("Origin")

// 		next.ServeHTTP(w, r)
// 	})
// }
