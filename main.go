package main

import (
	"log"
	"net/http"
	"os"

	"projet-forum/client/app"
)

func main() {
	application := app.InitApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	log.Printf("Serveur web lancé sur %s", address)
	if err := http.ListenAndServe(address, application.Router); err != nil {
		log.Fatalf("Erreur lancement serveur web : %v", err)
	}
}
