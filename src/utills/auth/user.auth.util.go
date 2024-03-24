package auth

import (
	"errors"
	"petpal-backend/src/models"
	user_utills "petpal-backend/src/utills/user"
	misc_utills "petpal-backend/src/utills/miscellaneous"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewUser(createUser models.CreateUser) (*models.User, error) {
	objID := primitive.NewObjectID().Hex()

	randomProfileImage, err := misc_utills.RandomProfileImage()
	if err != nil {
		return nil, err
	}

	// You can add more validation rules as needed
	newUser := &models.User{
		Individual: models.Individual{
			IndividualID: objID,
		},
		ID:                   objID,
		Username:             createUser.Username,
		Password:             createUser.Password,
		Email:                createUser.Email,
		FullName:             createUser.FullName,
		Address:              createUser.Address,
		DateOfBirth:          createUser.DateOfBirth,
		PhoneNumber:          createUser.PhoneNumber,
		ProfilePicture:       randomProfileImage,
		DefaultAccountNumber: "Deflut",
		DefaultBank:          "Deflut",
		Pets:                 []models.Pet{},
	}

	return newUser, nil
}

// RegisterHandler handles user registration
func RegisterUser(createUser models.CreateUser, db *models.MongoDB) (string, error) {
	// search if email already exists
	user, err := user_utills.GetUserByEmail(db, createUser.Email)
	if user != nil {
		return "", errors.New("User email alreay exist")
	}
	if err == nil {
		return "", errors.New("User email alreay exist")
	}

	// Hash the password securely
	hashedPassword, err := HashPassword(createUser.Password)
	if err != nil {
		return "", err
	}

	// Create a new user instance
	createUser.Password = hashedPassword
	newUser, err := NewUser(createUser)
	if err != nil {
		return "", err
	}

	// Insert the new user into the database
	newUser, err = user_utills.InsertUser(db, newUser)
	if err != nil {
		return "", err
	}

	// Generate a JWT token
	tokenString, err := GenerateToken(newUser.Username, newUser.Password, "user")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
