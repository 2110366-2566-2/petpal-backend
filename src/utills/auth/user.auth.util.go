package auth

import (
	"petpal-backend/src/models"
	user_utills "petpal-backend/src/utills/user"
)

func nextUserId() int {
	id := 5
	return id
}

func NewUser(createUser models.CreateUser) (*models.User, error) {
	newID := nextUserId()
	// You can add more validation rules as needed
	newUser := &models.User{
		Individual: models.Individual{
			IndividualID: newID,
		},
		Username:             createUser.Username,
		Password:             createUser.Password,
		Email:                createUser.Email,
		FullName:             createUser.FullName,
		Address:              createUser.Address,
		DateOfBirth:          createUser.DateOfBirth,
		PhoneNumber:          createUser.PhoneNumber,
		ProfilePicture:       []byte("Default"),
		DefaultAccountNumber: "Deflut",
		DefaultBank:          "Deflut",
		Pets:                 []models.Pet{},
	}

	return newUser, nil
}

// RegisterHandler handles user registration
func RegisterUser(createUser models.CreateUser, db *models.MongoDB) (string, error) {

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
