package utils

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, map[string]string{"error": message})
}
