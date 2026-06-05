package api

import (
	"net/http"
)

func main() {
	app := appInitApp()
	defer app.Close()

	logPrintf("Serveur lancé sur le port %d", app.Config.Port)
	serverErr := http.ListenAndServe(":8080", app.Router)
	if serverErr != nil {
		logFatal("Erreur lors du démarrage du serveur: %v", serverErr)
	}
}
