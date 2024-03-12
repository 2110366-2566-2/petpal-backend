package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"petpal-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"context"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ChangePassword(email string, newPassword string, login_type string, db *models.MongoDB) (string, error) {
	if login_type != "user" && login_type != "svcp" {
		return "Invalid login type", nil
	}
	var passwordKey string
	var emailKey string
	if login_type == "user" {
		passwordKey = "password"
		emailKey = "email"
	} else {
		passwordKey = "SVCPPassword"
		emailKey = "SVCPEmail"
	}
	// get collection
	collection := db.Collection(login_type)

	// find user by id
	var svcp models.User = models.User{}
	filter := bson.D{{Key: emailKey, Value: email}}
	err := collection.FindOne(context.Background(), filter).Decode(&svcp)
	if err != nil {
		return "SVCP/User not found (email=" + email + ")", err
	}

	// update user with new default bank account
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: passwordKey, Value: newPassword},
		}},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Failed to update SVCP/User password", err
	}

	return "", nil
}
