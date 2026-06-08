package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var ListTemp *template.Template

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {

	var buffer bytes.Buffer

	errRender := ListTemp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		log.Printf("Template error '%s': %v", name, errRender)
		http.Redirect(
			w,
			r,
			fmt.Sprintf(
				"/error?code=%d&message=%s",
				http.StatusInternalServerError,
				url.QueryEscape("Error loading the page"),
			),
			http.StatusSeeOther,
		)
		return
	}
	_, _ = buffer.WriteTo(w)
}
