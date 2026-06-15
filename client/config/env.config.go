// Package config charge les variables d'environnement et fournit des helpers de configuration.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv tente de charger un fichier .env à la racine du projet.
func LoadEnv() {
	errLoad := godotenv.Load("./.env")
	if errLoad != nil {
		log.Println("Aucun fichier .env trouve, utilisation des variables d'environnement systeme")
	}
}

// GetEnvWithDefault récupère une variable d'environnement ou renvoie une valeur par défaut.
func GetEnvWithDefault(key, defaultValue string) string {
	envVar, envErr := os.LookupEnv(key)
	if !envErr {
		return defaultValue
	}
	return envVar
}

// GetRequiredEnv récupère une variable d'environnement obligatoire ou arrête le programme en cas d'absence.
func GetRequiredEnv(key string) string {
	envVar, envErr := os.LookupEnv(key)
	if !envErr || envVar == "" {
		log.Fatalf("Erreur configuration - Variable d'environnement manquante : %s", key)
	}
	return envVar
}
