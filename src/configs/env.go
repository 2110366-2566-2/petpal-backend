package configs

import (
	"errors"
	"os"

	"github.com/lpernett/godotenv"
)

type EnvormentVariable struct {
	name                  string
	port                  string
	db_uri                string
	email_sender_address  string
	email_sender_password string
	jwt_secret            string
}

var instance *EnvormentVariable

func GetInstance() *EnvormentVariable {
	if instance == nil {
		instance = &EnvormentVariable{name: "Golang Singleton"}
	}
	return instance
}

func (s *EnvormentVariable) SetProductionEnv() error {

	err := godotenv.Load()
	if err != nil {
		return errors.New("Error loading .env file")
	}
	s.name = "Production"
	s.port = os.Getenv("PORT")
	s.db_uri = os.Getenv("DB_URI")
	s.email_sender_address = os.Getenv("EMAIL_SENDER_ADDRESS")
	s.email_sender_password = os.Getenv("EMAIL_SENDER_PASSWORD")
	s.jwt_secret = os.Getenv("JWT_SECRET")
	return nil
}
func (s *EnvormentVariable) SetTestEnv() error {
	err := godotenv.Load(".env.test")
	if err != nil {
		return errors.New("Error loading .env.test file")
	}
	s.name = "Test"
	s.port = os.Getenv("PORT")
	s.db_uri = os.Getenv("DB_URI")
	s.email_sender_address = os.Getenv("EMAIL_SENDER_ADDRESS")
	s.email_sender_password = os.Getenv("EMAIL_SENDER_PASSWORD")
	s.jwt_secret = os.Getenv("JWT_SECRET")
	return nil
}

func (s *EnvormentVariable) GetName() string {
	return s.name
}
func (s *EnvormentVariable) GetPort() string {
	return s.port
}
func (s *EnvormentVariable) GetDB_URI() string {
	return s.db_uri
}
func (s *EnvormentVariable) GetEmailSenderAddress() string {
	return s.email_sender_address
}
func (s *EnvormentVariable) GetEmailSenderPassword() string {
	return s.email_sender_password
}
func (s *EnvormentVariable) GetJWT_SECRET() string {
	return s.jwt_secret
}
