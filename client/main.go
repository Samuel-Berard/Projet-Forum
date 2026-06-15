// Package main démarre l'application cliente (front web) et configure le serveur HTTP.
package main

import (
	"log"
	"net/http"

	"projet-forum/client/app"
)

// main initialise l'application cliente et démarre le serveur HTTP sur le port 3000.
func main() {
	app := app.InitApp()

	log.Printf("Serveur Frontend (Client) lancé : http://localhost:3000")
	serveErr := http.ListenAndServe(":3000", app.Router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur client - %s", serveErr.Error())
	}
}
