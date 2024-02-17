package auth

import (
	"petpal-backend/src/models"

	"errors"

	"github.com/gin-gonic/gin"
)

type CurrentEntity interface {
	// Define methods shared by both models
}

func GetCurrentUserByGinContext(c *gin.Context, db *models.MongoDB) (*models.User, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	// Parse request body to get user data
	entity, err := GetCurrentEntity(token, db)
	if err != nil {
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.User:
		return entity, nil
		// Handle user
	case *models.SVCP:
		return nil, errors.New(" Need token of type User but recives SVCP type")
		// Handle svcp
	}
	return nil, errors.New(" Need token of type User but wrong type")
}

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
		ProfilePicture:       "Deflut",
		DefaultAccountNumber: "Deflut",
		DefaultBank:          "Deflut",
		Pets:                 nil,
	}

	return newUser, nil
}
