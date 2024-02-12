package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertUser(db *models.MongoDB, user *models.User) (*models.User, error) {
	// Get the users collection
	collection := db.Collection("user")

	// Insert the user into the collection
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	// Return the inserted user
	return user, nil
}

func GetUserByEmail(db *models.MongoDB, email string) (*models.User, error) {
	// get collection
	collection := db.Collection("user")
	// find user by email
	// note: IndividualID is not present in the database yet, so this always returns an error not found
	var user models.User = models.User{}
	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserPet(db *models.MongoDB, userEmail string) (*[]models.Pet, error) {
	user, err := GetUserByEmail(db, userEmail)
	if err != nil {
		return nil, err
	}

	if user.Pets == nil {
		emptySlice := make([]models.Pet, 0)
		return &emptySlice, nil
	}
	// add pets ownername
	for i := range user.Pets {
		user.Pets[i].OwnerUsername = user.Username
	}
	return &user.Pets, nil
}

func SetDefaultBankAccount(email string, defaultAccountNumber string, defaultBank string, db *models.MongoDB) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by id
	var user models.User = models.User{}
	filter := bson.D{{Key: "email", Value: email}}
	err := user_collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "User not found (email=" + email + ")", err
	}

	// update user with new default bank account
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "defaultAccountNumber", Value: defaultAccountNumber},
			{Key: "defaultBank", Value: defaultBank},
		}},
	}
	_, err = user_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Failed to update user", err
	}

	return "", nil
}

func ChangePassword(email string, newPassword string, db *models.MongoDB) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by id
	var user models.User = models.User{}
	filter := bson.D{{Key: "email", Value: email}}
	err := user_collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "User not found (email=" + email + ")", err
	}

	// update user with new default bank account
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: newPassword},
		}},
	}
	_, err = user_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Failed to update user password", err
	}

	return "", nil
}
