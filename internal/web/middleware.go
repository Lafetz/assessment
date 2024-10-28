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
