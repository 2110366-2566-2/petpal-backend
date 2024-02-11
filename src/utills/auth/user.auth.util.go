package auth

import (
	"errors"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"
	"time"
)

func GetCurrnetUser(token string, db *models.MongoDB) (*models.User, error) {
	loginRes, err := DecodeToken(token)
	if err != nil {
		return nil, err
	}
	loginType := loginRes.LoginType
	if loginType == "user" {
		user, err := utills.GetUserByEmail(db, loginRes.UserEmail)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		return nil, errors.New("Get Wrong User type we only accept svcp login type but get " + loginType)
	}
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
		Address:              "Defult",
		DateOfBirth:          time.Now(),
		PhoneNumber:          "Deflut",
		ProfilePicture:       "Deflut",
		DefaultAccountNumber: "Deflut",
		DefaultBank:          "Deflut",
		Pets:                 nil,
	}

	return newUser, nil
}
