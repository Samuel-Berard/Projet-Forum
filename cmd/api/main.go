package main

import (
	"log"
	"net/http"

	"projet-forum/src/api/config"
	"projet-forum/src/api/controllers"
	"projet-forum/src/api/repositories"
	"projet-forum/src/api/router"
	"projet-forum/src/api/services"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Port   string
}

func InitApp() *App {

	config.LoadEnv()

	db := config.InitDB()

	userRepo := repositories.NewUtilisateurRepository(db)
	filRepo := repositories.NewFilRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	userService := services.NewUtilisateurService(userRepo)
	filService := services.NewFilService(filRepo)
	messageService := services.NewMessageService(messageRepo, filRepo)

	userController := controllers.NewUtilisateurController(userService)
	filController := controllers.NewFilController(filService)
	messageController := controllers.NewMessageController(messageService)

	r := router.NewRouter(userController, filController, messageController)

	port := config.GetEnvWithDefault("PORT", "8080")

	log.Printf("Serveur lancé sur le port %s", port)

	return &App{
		Router: r,
		Port:   port,
	}
}

func main() {
	app := InitApp()

	serverErr := http.ListenAndServe(":"+app.Port, app.Router)
	if serverErr != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", serverErr)
	}
}
