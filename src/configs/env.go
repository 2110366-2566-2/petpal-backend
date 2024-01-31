package configs

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

func GetPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading port from .env file")
	}
	port := os.Getenv("PORT")
	return port
}

func GetDB_URI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading db_uri from .env file")
	}
	db_uri := os.Getenv("DB_URI")
	return db_uri
}
