package auth

import (
	"errors"
	"fmt"
	"petpal-backend/src/models"
	utills "petpal-backend/src/utills/serviceprovider"
)

func GetCurrentSVCP(token string, db *models.MongoDB) (*models.SVCP, error) {
	loginRes, err := DecodeToken(token)
	if err != nil {
		return nil, err
	}
	loginType := loginRes.LoginType
	if loginType == "svcp" {
		user, err := utills.GetSVCPByEmail(db, loginRes.UserEmail)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		return nil, errors.New("Get Wrong User type we only accept user login type but get " + loginType)
	}
}

func nextSVCPId() int {
	id := 5
	return id
}

func NewSVCP(createSVCP models.CreateSVCP) (*models.SVCP, error) {
	newID := nextUserId()
	fmt.Println(createSVCP)
	// You can add more validation rules as needed
	newSVCP := &models.SVCP{
		Individual: models.Individual{
			IndividualID: newID,
		},
		SVCPID:                "Defult",
		SVCPImg:               "Defult",
		SVCPUsername:          createSVCP.SVCPUsername,
		SVCPPassword:          createSVCP.SVCPPassword,
		SVCPEmail:             createSVCP.SVCPEmail,
		IsVerified:            false,
		SVCPResponsiblePerson: "Defult",
		DefaultBank:           "Defult",
		DefaultAccountNumber:  "Defult",
		License:               "Defult",
		Location:              "Defult",
		SVCPAdditionalImg:     "Defult",
		SVCPServiceType:       createSVCP.SVCPServiceType,
	}

	return newSVCP, nil
}