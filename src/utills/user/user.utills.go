package utills

import (
	"context"
	"petpal-backend/src/models"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetUsers(db *models.MongoDB, filter bson.D, page int64, per int64) ([]models.User, error) {
	collection := db.Collection("user")

	// define options for pagination
	opts := options.Find().SetSkip(page * per).SetLimit(per)

	// Find all documents in the collection
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, err
}

func GetUserByID(db *models.MongoDB, id string) (*models.User, error) {
	// get collection
	collection := db.Collection("user")

	// find user by id
	var user models.User = models.User{}
	filter := bson.D{{Key: "_id", Value: id}}
	opts := options.FindOne().SetProjection(bson.D{{Key: "search_history", Value: 0}})
	err := collection.FindOne(context.Background(), filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(db *models.MongoDB, email string) (*models.User, error) {
	// get collection
	collection := db.Collection("user")
	// find user by email
	// note: IndividualID is not present in the database yet, so this always returns an error not found
	var user models.User = models.User{}
	filter := bson.D{{Key: "email", Value: email}}
	opts := options.FindOne().SetProjection(bson.D{{Key: "search_history", Value: 0}})
	err := collection.FindOne(context.Background(), filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetSearchHistory(db *models.MongoDB, id string) ([]models.SearchHistory, error) {
	// get collection
	collection := db.Collection("user")

	var search_history models.UserSearchHistory = models.UserSearchHistory{}
	filter := bson.D{{Key: "_id", Value: id}}
	opts := options.FindOne().SetProjection(bson.D{{Key: "search_history", Value: 1}})
	_ = collection.FindOne(context.Background(), filter, opts).Decode(&search_history)

	return search_history.SearchHistory, nil
}

func AddSearchHistory(db *models.MongoDB, id string, search_filters models.SearchFilter) error {
	// get collection
	collection := db.Collection("user")

	// find user by id
	filter := bson.D{{Key: "_id", Value: id}}

	search_history := models.SearchHistory{
		Date:         time.Now(),
		SearchFilter: search_filters,
	}

	// update user
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "search_history", Value: search_history},
		}},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(db *models.MongoDB, userUpdate *bson.M, userId string) (string, error) {
	// get collection
	collection := db.Collection("user")

	// find user by id
	user, err := GetUserByID(db, userId)
	if err != nil {
		return "User not found", err
	}

	for key, value := range *userUpdate {
		user.UpdateField(key, value)
	}
	filter := bson.D{{Key: "_id", Value: userId}}
	// update user
	update := bson.D{
		{Key: "$set", Value: user},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Failed to update user", err
	}

	return "", nil

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

func DeleteBankAccount(email string, db *models.MongoDB) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by id
	var user models.User = models.User{}
	filter := bson.D{{Key: "email", Value: email}}
	err := user_collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "User not found (email=" + email + ")", err
	}

	// update default bank account to empty
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "defaultAccountNumber", Value: ""},
			{Key: "defaultBank", Value: ""},
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

		type user_decode struct {
			Email          string `json:"email" bson:"email"`
			ProfilePicture []byte `json:"profilePicture" bson:"profilePicture"`
		}

		// Find the user by email in the "user" collection
		var user user_decode = user_decode{}
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
			SVCPImg   []byte `json:"SVCPImg" bson:"SVCPImg"`
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
