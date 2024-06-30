package common

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any, status int) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func WriteError(w http.ResponseWriter, message string, status int) error {
	return WriteJSON(w, map[string]any{"error": message}, status)
}
