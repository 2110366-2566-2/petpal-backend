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

// UploadProfileImage uploads a profile image for a user and svcp.
// It takes the email, file content (image bytes), and a MongoDB instance.
// Returns a gin.H (response) and an error if any.
func UploadProfileImage(email string, fileContent []byte, userType string, db *models.MongoDB) (gin.H, error) {

	if userType == "user" {
		userCollection := db.Collection("user")

		// Find the user by email in the "user" collection
		var user models.User = models.User{}
		filter := bson.D{{Key: "email", Value: email}}
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
		return gin.H{"message": "Profile image stored successfully in 'user' collection", "updated": results}, nil
	} else if userType == "svcp" {
		userCollection := db.Collection("svcp")

		//temporarily struct for svcp
		type Svcp struct {
			SVCPEmail string `json:"SVCPEmail"`
			SVCPImg   string `json:"SVCPImg"`
		}

		// Find the svcp by email in the "svcp" collection
		var user Svcp = Svcp{}

		filter := bson.D{{Key: "SVCPEmail", Value: email}}
		err := userCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			// If an error occurs during the database query, return an error response
			return gin.H{"error": err.Error() + email}, err
		}

		// Update the svcp's profilePicture field with the new file content
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "SVCPImg", Value: fileContent},
			}},
		}

		// Perform the update operation in the "svcp" collection
		results, err := userCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			// If an error occurs during the update, return an error response
			return gin.H{"error": "Error updating profile Picture in the database"}, err
		}

		// Return a success message along with the update results
		return gin.H{"message": "Profile image stored successfully in 'svcp' collection", "updated": results}, nil
	}
	return gin.H{"error": "missing usertype in backend"}, nil
}

func GetProfileImage(email string, userType string, db *models.MongoDB) (gin.H, error) {

	if userType == "user" {
		// Access the "user" collection in the MongoDB database
		userCollection := db.Collection("user")

		// Find the user by email in the "user" collection
		var user models.User = models.User{}
		filter := bson.D{{Key: "email", Value: email}}
		err := userCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			// If an error occurs during the database query, return an error response
			return gin.H{"error": err.Error()}, err
		}

		// Return the profile picture file content
		results := user.ProfilePicture
		return gin.H{"message": "Get profile image from 'user' collection", "email": email, "result": results}, nil
	} else if userType == "svcp" {

		//temporarily struct for svcp
		type Svcp struct {
			SVCPEmail string `json:"SVCPEmail"`
			SVCPImg   string `json:"SVCPImg"`
		}

		// Access the "svcp" collection in the MongoDB database
		userCollection := db.Collection("svcp")

		// Find the user by email in the "svcp" collection
		var user Svcp = Svcp{}
		filter := bson.D{{Key: "SVCPEmail", Value: email}}
		err := userCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			// If an error occurs during the database query, return an error response
			return gin.H{"error": err.Error()}, err
		}

		// Return the profile picture file content
		results := user.SVCPImg
		return gin.H{"message": "Get profile image from 'svcp' collection", "email": email, "result": results}, nil
	}
	return gin.H{"error": "missing usertype in backend"}, nil
}
