package httpx

import (
	"encoding/json"
	"net/http"
)

// JSON sends a JSON response with the given status code.
func JSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

// OK sends a 200 OK JSON response.
func OK(
	w http.ResponseWriter,
	data any,
) {
	JSON(w, http.StatusOK, data)
}

// Error sends a JSON error response.
func Error(
	w http.ResponseWriter,
	status int,
	message string,
) {
	JSON(w, status, map[string]string{
		"error": message,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
