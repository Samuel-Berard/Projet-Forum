package app

import (
	"database/sql"

	"projet-forum/api/config"
	"projet-forum/api/controllers"
	"projet-forum/api/repositories"
	"projet-forum/api/routers"
	"projet-forum/api/services"

	"github.com/gorilla/mux"
)

type App struct {
	Db     *sql.DB
	Router *mux.Router
}

func InitApp() *App {
	// Chargement des variables d'environnement
	config.LoadEnv()

	// Initialisation de la connexion à la base de données
	db := config.InitDB()

	// Initialisation des repositories
	utilisateurRepository := repositories.InitUtilisateurRepository(db)
	filRepository := repositories.InitFilRepository(db)
	messageRepository := repositories.InitMessageRepository(db)

	// Initialisation des services
	utilisateurService := services.InitUtilisateurService(utilisateurRepository)
	filService := services.InitFilService(filRepository)
	messageService := services.InitMessageService(messageRepository, filRepository)
	actusService := services.InitActusService()

	// Initialisation des controllers
	utilisateurController := controllers.InitUtilisateurController(utilisateurService)
	filController := controllers.InitFilController(filService)
	messageController := controllers.InitMessageController(messageService)
	actusController := controllers.InitActusController(actusService)

	// Enregistrement des routes (avec ajout du préfixe "/api/...")
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	routers.RegisterUtilisateurRoutes(router, utilisateurController)
	routers.RegisterFilRoutes(router, filController)
	routers.RegisterMessageRoutes(router, messageController)
	routers.RegisterActusRoutes(router, actusController)

	return &App{
		Db:     db,
		Router: router,
	}
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
}
