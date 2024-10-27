package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	person "github.com/lafetz/assessment/internal/core/service"
	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		expectedSize int
		expectedPage int
	}{
		{
			name:         "No parameters",
			queryParams:  "",
			expectedSize: 10,
			expectedPage: 1,
		},
		{
			name:         "Valid parameters",
			queryParams:  "?size=20&page=2",
			expectedSize: 20,
			expectedPage: 2,
		},
		{
			name:         "Invalid size",
			queryParams:  "?size=abc&page=2",
			expectedSize: 10,
			expectedPage: 2,
		},
		{
			name:         "Negative size",
			queryParams:  "?size=-5&page=2",
			expectedSize: 10,
			expectedPage: 2,
		},
		{
			name:         "Valid size and invalid page",
			queryParams:  "?size=20&page=0",
			expectedSize: 20,
			expectedPage: 1,
		},
		{
			name:         "Valid size and page with string",
			queryParams:  "?size=10&page=abc",
			expectedSize: 10,
			expectedPage: 1,
		},
		{
			name:         "Page out of bounds",
			queryParams:  "?size=10&page=5",
			expectedSize: 10,
			expectedPage: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/example"+tt.queryParams, nil)
			pagination := ParsePagination(req)

			assert.Equal(t, tt.expectedSize, pagination.Size)
			assert.Equal(t, tt.expectedPage, pagination.Page)
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Not Found Error",
			err:          person.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "not found",
		},
		{
			name:         "Generic Error",
			err:          errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "internal server error",
		},
		{
			name:         "Nil Error",
			err:          nil,
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleError(tt.err, w, slog.Default())
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusInternalServerError || tt.err != nil {
				var response errorMessage
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, response.Message)
				assert.Equal(t, tt.expectedCode, response.StatusCode)
			}

		})
	}
}
