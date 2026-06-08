package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSON envoie une réponse JSON avec le code HTTP donné
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

// WriteError envoie une réponse d'erreur JSON
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// WriteLog affiche un message dans les logs du serveur
func WriteLog(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
