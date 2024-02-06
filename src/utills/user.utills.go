package utills

import (
	"petpal-backend/src/models"
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
		CreateUser:           createUser,
		PhoneNumber:          "Mock",
		ProfilePicture:       "Mock",
		DefaultAccountNumber: "Mock",
		DefaultBank:          "Mock",
		Pet:                  nil,
	}

	return newUser, nil
}
