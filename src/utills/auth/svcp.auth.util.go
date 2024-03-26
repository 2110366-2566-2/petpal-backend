package auth

import (
	"errors"
	"petpal-backend/src/models"
	misc_utills "petpal-backend/src/utills/miscellaneous"
	svcp_utills "petpal-backend/src/utills/serviceprovider"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewSVCP(createSVCP models.CreateSVCP) (*models.SVCP, error) {
	objID := primitive.NewObjectID().Hex()

	randomProfileImage, err := misc_utills.RandomProfileImage()
	if err != nil {
		return nil, err
	}
	randomServiceImage, err2 := misc_utills.RandomServiceImage()
	if err2 != nil {
		return nil, err2
	}

	// You can add more validation rules as needed
	newSVCP := &models.SVCP{
		Individual: models.Individual{
			IndividualID: objID,
		},
		SVCPID:                objID,
		SVCPImg:               randomProfileImage,
		SVCPUsername:          createSVCP.SVCPUsername,
		SVCPPassword:          createSVCP.SVCPPassword,
		SVCPEmail:             createSVCP.SVCPEmail,
		IsVerified:            false,
		SVCPResponsiblePerson: "",
		DefaultBank:           "",
		DefaultAccountNumber:  "",
		License:               "",
		Address:               "",
		SVCPAdditionalImg:     randomServiceImage,
		SVCPServiceType:       createSVCP.SVCPServiceType,
		Services:              []models.Service{},
	}

	return newSVCP, nil
}

// RegisterHandler handles user registration
func RegisterSVCP(createSVCP models.CreateSVCP, db *models.MongoDB) (string, error) {

	// search if email already exists
	svcp, err := svcp_utills.GetSVCPByEmail(db, createSVCP.SVCPEmail)
	if svcp != nil {
		return "", errors.New("SVCP email alreay exist")
	}
	if err == nil {
		return "", errors.New("SVCP email alreay exist")
	}

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
