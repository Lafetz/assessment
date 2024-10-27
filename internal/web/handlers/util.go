package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	person "github.com/lafetz/assessment/internal/core/service"
)

type PaginationParams struct {
	Size int
	Page int
}

func ParsePagination(r *http.Request) PaginationParams {
	sizeStr := r.URL.Query().Get("size")
	pageStr := r.URL.Query().Get("page")

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	return PaginationParams{
		Size: size,
		Page: page,
	}
}
func HandleError(err error, w http.ResponseWriter, logger *slog.Logger) {

	if err != nil {
		switch {
		case errors.Is(err, person.ErrNotFound):
			writeError(w, "not found", http.StatusNotFound)
		default:
			logger.Error(err.Error())
			writeError(w, "internal server error", http.StatusInternalServerError)

		}
	} else {
		logger.Error("expected error but got nil")
		writeError(w, "internal server error", http.StatusInternalServerError)
	}
}

type errorMessage struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorMessage{
		StatusCode: statusCode,
		Message:    message,
	})
}
