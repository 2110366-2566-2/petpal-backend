package auth

import (
	"errors"
	"petpal-backend/src/models"
	user_utills "petpal-backend/src/utills/user"
)

func GetCurrentUser(token string, db *models.MongoDB) (*models.User, error) {
	loginRes, err := DecodeToken(token)
	if err != nil {
		return nil, err
	}
	loginType := loginRes.LoginType
	if loginType == "user" {
		user, err := user_utills.GetUserByEmail(db, loginRes.UserEmail)
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
