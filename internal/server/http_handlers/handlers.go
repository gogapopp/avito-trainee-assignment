package httphandlers

import (
	"assignment/internal/models"
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func errorJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	jsonResponse(w, statusCode, models.ErrorResponse{Message: message})
}
