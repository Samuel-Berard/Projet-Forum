package main

import (
	"log"
	"net/http"

	"projet-forum/src/client/api"
	"projet-forum/src/client/controllers"
	"projet-forum/src/client/router"
	"projet-forum/src/client/services"
)

type ClientApp struct {
	Router http.Handler
	Port   string
}

func InitClientApp() *ClientApp {

	apiBaseURL := "http://localhost:8080/api"

	forumApi := api.NewForumApi(apiBaseURL)
	forumService := services.NewForumClientService(forumApi)
	viewController := controllers.NewViewController(forumService)

	r := router.NewClientRouter(viewController)

	port := "3000"

	log.Printf("Serveur Frontend (Client) lancé sur http://localhost:%s", port)

	return &ClientApp{
		Router: r,
		Port:   port,
	}
}

func main() {
	app := InitClientApp()

	serverErr := http.ListenAndServe(":"+app.Port, app.Router)
	if serverErr != nil {
		log.Fatalf("Erreur lors du démarrage du serveur client: %v", serverErr)
	}
}
