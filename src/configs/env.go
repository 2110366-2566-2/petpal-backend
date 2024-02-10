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

func GetEmailSenderAddress() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading email_sender_address from .env file")
	}
	email_sender_address := os.Getenv("EMAIL_SENDER_ADDRESS")
	return email_sender_address
}

func GetEmailSenderPassword() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading email_sender_passwordfrom .env file")
	}
	email_sender_password := os.Getenv("EMAIL_SENDER_PASSWORD")
	return email_sender_password
}
