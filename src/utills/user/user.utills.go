package utills

import (
	"context"
	"fmt"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		Username:             createUser.Username,
		Password:             createUser.Password,
		Email:                createUser.Email,
		FullName:             createUser.FullName,
		PhoneNumber:          "Mock",
		ProfilePicture:       "Mock",
		DefaultAccountNumber: "Mock",
		DefaultBank:          "Mock",
		Pets:                 nil,
	}

	return newUser, nil
}

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
	opts := options.Find().SetSkip(page * per).SetLimit(per)

	// Find all documents in the collection
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode results
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user models.User = models.User{}
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
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

func AddUserPet(db *models.MongoDB, pet *models.Pet, user_id string) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by email
	user_objectid, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return "Invalid user id", err
	}

	filter := bson.D{{Key: "_id", Value: user_objectid}}
	res, err := user_collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$push", Value: bson.D{{Key: "pets", Value: pet}}}})
	if res.MatchedCount == 0 {
		return "User not found (id=" + user_id + ")", err
	}
	if err != nil {
		return "Failed to add pet", err
	}

	return "", nil
}

func UpdateUserPet(db *models.MongoDB, pet *models.Pet, user_id string, pet_idx int) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by email
	user_objectid, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return "Invalid user id", err
	}
	filter := bson.D{{Key: "_id", Value: user_objectid}}
	res, err := user_collection.UpdateOne(context.Background(), filter, bson.D{{
		Key: "$set", Value: bson.D{{
			Key: "pets." + fmt.Sprint(pet_idx), Value: pet,
		}},
	}})
	if res.MatchedCount == 0 {
		return "User not found (id=" + user_id + ")", err
	}
	if err != nil {
		return "Failed to update pet", err
	}

	return "", nil

}

// deletes a pet from a user's pet list by index
// note that when index is out of range, it will do nothing and *NOT* return an error
func DeleteUserPet(db *models.MongoDB, user_id string, pet_idx int) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by email
	user_objectid, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return "Invalid user id", err
	}

	filter := bson.D{{Key: "_id", Value: user_objectid}}
	res, err := user_collection.UpdateOne(context.Background(), filter, bson.A{bson.D{{
		Key: "$set", Value: bson.D{{
			Key: "pets", Value: bson.D{{
				Key: "$concatArrays", Value: bson.A{
					bson.D{{Key: "$slice", Value: bson.A{"$pets", pet_idx}}},
					bson.D{{Key: "$slice", Value: bson.A{"$pets", pet_idx + 1, bson.D{{Key: "$size", Value: "$pets"}}}}},
				},
			}},
		}},
	}}})

	if res.MatchedCount == 0 {
		return "User not found (id=" + user_id + ")", err
	}

	if err != nil {
		return "Failed to delete pet", err
	}

	return "", nil
}

func UpdateUser(db *models.MongoDB, user *bson.M, id string) (string, error) {
	// get collection
	collection := db.Collection("user")

	// find user by id
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "Invalid user id", err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}

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
