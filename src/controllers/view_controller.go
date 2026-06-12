package controllers

import (
	"html/template"
	"log"
	"net/http"
)

func AfficherAccueil(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {

		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		log.Println("Erreur Template :", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
