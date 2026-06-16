package controllers

import (
	"fmt"
	"log"
	"net/http"
)

type CookiesControllers struct {
}

func InitCookiesController() *CookiesControllers {
	return &CookiesControllers{}
}

func AddCoockieHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "key_cookies",
		Value:    "valeur du cookies",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   15 * 60,
	})

	fmt.Fprintf(w, "Cookie ajouté avec la clé : key_cookies")
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, cookieErr := r.Cookie("key_cookies")
	if cookieErr != nil {
		log.Println(cookieErr)
		http.Error(w, "Erreur lecture du cookie", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Valeur du cookie : %v", cookie.Value)
}

func DeleteCookieHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "key_cookies",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	fmt.Fprintf(w, "Cookie spprimé avec la clé : key_cookies")
}
