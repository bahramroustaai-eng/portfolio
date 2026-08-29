package transport

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"portfolio/internal/apperr"
)

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	writeJSON(w, statusCode, ErrorResponse{Error: msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		writeError(w, statusForCode(appErr.Code), err.Error())
		return
	}
	slog.Error("unhandled error", "error", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

func statusForCode(c apperr.Code) int {
	switch c {
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeConflict:
		return http.StatusConflict
	case apperr.CodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}
