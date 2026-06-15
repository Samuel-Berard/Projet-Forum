// Package templates gère le chargement et le rendu des fichiers HTML de l'application cliente.
package templates

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

// TemplateManager contient l'ensemble des templates chargés en mémoire.
type TemplateManager struct {
	templates *template.Template
}

// NewTemplatesManager construit un TemplateManager et charge les templates.
func NewTemplatesManager() *TemplateManager {
	tm := &TemplateManager{}
	tm.load()
	return tm
}

// load charge tous les fichiers HTML depuis le dossier ./client/templates.
func (tm *TemplateManager) load() {
	// Fonctions utilitaires disponibles dans les templates.
	funcs := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	loadedTemplates, loadedTemplatesErr := template.New("").Funcs(funcs).ParseGlob("./client/templates/*.html")
	if loadedTemplatesErr != nil {
		log.Fatalf("Erreur template - %s", loadedTemplatesErr.Error())
		return
	}
	tm.templates = loadedTemplates
	log.Printf("Template - Templates chargés : %s", tm.templates.DefinedTemplates())
}

// RenderTemplate exécute le template demandé et envoie le résultat HTTP au client.
func (tm *TemplateManager) RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	var buffer bytes.Buffer
	if tm.templates == nil {
		http.Error(w, "Templates non initialises", http.StatusInternalServerError)
		return
	}

	errRender := tm.templates.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		log.Printf("Erreur template '%s' : %v", name, errRender)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	buffer.WriteTo(w)
}
