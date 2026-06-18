// Package app assemble les composants de l'application cliente et configure le routeur.
package app

import (
	"projet-forum/client/api"
	"projet-forum/client/config"
	"projet-forum/client/controllers"
	"projet-forum/client/routers"
	"projet-forum/client/services"
	"projet-forum/client/templates"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

// InitApp initialise l'application, charge les templates et configure le routeur.
func InitApp() *App {
	config.LoadEnv()

	templatesManager := templates.NewTemplatesManager()

	// URL de base de l'API consommée par le client.
	baseURL := config.GetEnvWithDefault("API_BASE_URL", "http://localhost:8080/api")
	forumApi := api.InitForumApi(baseURL)

	forumService := services.InitForumService(forumApi)

	forumController := controllers.InitForumController(forumService, templatesManager)

	router := mux.NewRouter()
	routers.RegisterAssetsRoutes(router)
	routers.RegisterForumRoutes(router, forumController)

	return &App{
		Router: router,
	}
}
