package auth

import (
	"errors"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"
)

func GetCurrnetUser(token string, db *models.MongoDB) (*models.User, error) {
	loginRes, err := DecodeToken(token)
	if err != nil {
		return nil, err
	}
	loginType := loginRes.LoginType
	if loginType == "user" {
		user, err := utills.GetUserByEmail(db, loginRes.Username)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		return nil, errors.New("Get Wrong User type we only accept user login type but get " + loginType)
	}
}
