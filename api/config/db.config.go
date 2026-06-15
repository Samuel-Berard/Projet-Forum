package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	// Recuperation des parametres lies a la base de donnees.
	user := GetRequiredEnv("DB_USER")
	pwd := GetEnvWithDefault("DB_PWD", "")
	host := GetRequiredEnv("DB_HOST")
	port := GetRequiredEnv("DB_PORT")
	name := GetRequiredEnv("DB_NAME")

	// Preparation de la chaine de connexion a la base de donnees.
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, name)

	// Mise en place de la connexion.
	dbContext, dbContextErr := sql.Open("mysql", connectionString)
	if dbContextErr != nil {
		log.Fatalf("Erreur connection base de donnees - Erreur : \n\t %s", dbContextErr.Error())
	}

	// Test de ping la base de donnees.
	pingErr := dbContext.Ping()
	if pingErr != nil {
		dbContext.Close()
		log.Fatalf("Erreur ping base de donnees - Erreur : \n\t %s", pingErr.Error())
	}

	log.Printf("BDD - Connexion reussie")
	return dbContext
}
