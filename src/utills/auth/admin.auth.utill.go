package auth

import (
	"errors"
	"petpal-backend/src/models"
	admin_utills "petpal-backend/src/utills/admin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewAdmin(createUser models.CreateAdmin) (*models.Admin, error) {
	objID := primitive.NewObjectID().Hex()
	// You can add more validation rules as needed
	newAdmin := &models.Admin{
		Individual: models.Individual{
			IndividualID: objID,
		},
		AdminID: objID,
		Email:   createUser.Email,
		FullName: createUser.FullName,
		Password: createUser.Password,
		Username: createUser.Username,
	}

	return newAdmin, nil
}

// RegisterHandler handles user registration
func RegisterAdmin(createAdmin models.CreateAdmin, db *models.MongoDB) (string, error) {
	// search if email already exists
	user, err := admin_utills.GetAdminByEmail(db, createAdmin.Email)
	if user != nil {
		return "", errors.New("admin email alreay exist")
	}
	if err == nil {
		return "", errors.New("admin email alreay exist")
	}

	// Hash the password securely
	hashedPassword, err := HashPassword(createAdmin.Password)
	if err != nil {
		return "", err
	}

	// Create a new user instance
	createAdmin.Password = hashedPassword
	newUser, err := NewAdmin(createAdmin)
	if err != nil {
		return "", err
	}

	// Insert the new user into the database
	newUser, err = admin_utills.InsertAdmin(db, newUser)
	if err != nil {
		return "", err
	}

	// Generate a JWT token
	tokenString, err := GenerateToken(newUser.Username, newUser.Password, "admin")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}