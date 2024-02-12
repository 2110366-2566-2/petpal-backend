package utills

import (
	"context"
	"petpal-backend/src/models"
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