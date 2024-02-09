package utills

import (
	"context"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
		Pets:                 nil,
	}

	return newUser, nil
}

func SetDefaultBankAccount(username string, defaultBankAccountNumber string, defaultBank string, db *models.MongoDB) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by id
	var user models.User = models.User{}
	filter := bson.D{{Key: "username", Value: username}}
	err := user_collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "User not found (" + username + ")", err
	}

	// update user with new default bank account
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "defaultAccountNumber", Value: defaultBankAccountNumber},
			{Key: "defaultBank", Value: defaultBank},
		}},
	}
	_, err = user_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Failed to update user", err
	}

	return "", nil
}

// UploadProfileImage uploads a profile image for a user.
// It takes the username, file content (image bytes), and a MongoDB instance.
// Returns a gin.H (response) and an error if any.
func UploadProfileImage(username string, fileContent []byte, db *models.MongoDB) (gin.H, error) {
	// Access the "user" collection in the MongoDB database
	userCollection := db.Collection("user")

	// Find the user by username in the "user" collection
	var user models.User = models.User{}
	filter := bson.D{{Key: "username", Value: username}}
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		// If an error occurs during the database query, return an error response
		return gin.H{"error": err.Error()}, err
	}

	// Update the user's profilePicture field with the new file content
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "profilePicture", Value: fileContent},
		}},
	}

	// Perform the update operation in the "user" collection
	results, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If an error occurs during the update, return an error response
		return gin.H{"error": "Error updating profile Picture in the database"}, err
	}

	// Return a success message along with the update results
	return gin.H{"message": "Profile Picture File stored successfully in 'user' collection", "updated": results}, nil
}

func GetProfileImage(username string, db *models.MongoDB) (string, error) {
	// Access the "user" collection in the MongoDB database
	userCollection := db.Collection("user")

	// Find the user by username in the "user" collection
	var user models.User = models.User{}
	filter := bson.D{{Key: "username", Value: username}}
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		// If an error occurs during the database query, return an error response
		return "", err
	}

	// Return the profile picture file content
	return user.ProfilePicture, nil
}
