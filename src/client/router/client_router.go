package router

import (
	"net/http"
	"projet-forum/src/client/controllers"

	"github.com/gorilla/mux"
)

func NewClientRouter(viewController *controllers.ViewController) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", viewController.AfficherAccueil).Methods("GET")
	r.HandleFunc("/login", viewController.AfficherLogin).Methods("GET")
	r.HandleFunc("/signup", viewController.AfficherSignup).Methods("GET")
	r.HandleFunc("/forum", viewController.AfficherForum).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return r
}
