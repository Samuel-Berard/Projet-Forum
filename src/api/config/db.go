package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {

	user := GetRequiredEnv("DB_USER")
	pwd := GetEnvWithDefault("DB_PWD", "")
	host := GetRequiredEnv("DB_HOST")
	port := GetRequiredEnv("DB_PORT")
	name := GetRequiredEnv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, name)

	dbContext, dbContextErr := sql.Open("mysql", connectionString)
	if dbContextErr != nil {
		log.Fatalf("Erreur connection base de donnees - Erreur : \n\t %s", dbContextErr.Error())
	}

	pingErr := dbContext.Ping()
	if pingErr != nil {
		dbContext.Close()
		log.Fatalf("Erreur ping base de donnees - Erreur : \n\t %s", pingErr.Error())
	}

	log.Printf("BDD - Connexion reussie")
	return dbContext
}
