package auth

import (
	"fmt"
	"petpal-backend/src/models"
	svcp_utills "petpal-backend/src/utills/serviceprovider"
)

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

// RegisterHandler handles user registration
func RegisterSVCP(createSVCP models.CreateSVCP, db *models.MongoDB) (string, error) {

	// Hash the password securely
	hashedPassword, err := HashPassword(createSVCP.SVCPPassword)
	if err != nil {
		return "", err
	}

	// Create a new user instance
	createSVCP.SVCPPassword = hashedPassword
	newSVCP, err := NewSVCP(createSVCP)
	if err != nil {
		return "", err
	}

	// Insert the new user into the database
	newSVCP, err = svcp_utills.InsertSVCP(db, newSVCP)
	if err != nil {
		return "", err
	}

	// Generate a JWT token
	tokenString, err := GenerateToken(newSVCP.SVCPUsername, newSVCP.SVCPEmail, "svcp")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
