package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(db *models.MongoDB, filter bson.D, page int64, per int64) ([]models.User, error) {
	collection := db.Collection("user")
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

func GetUserByID(db *models.MongoDB, id int64) (*models.User, error) {
	// get collection
	collection := db.Collection("user")

	// find user by id
	// note: IndividualID is not present in the database yet, so this always returns an error not found
	var user models.User = models.User{}
	filter := bson.D{{Key: "IndividualID", Value: id}}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

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
		Pets:                  nil,
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
		return "User not found ("+username+")", err
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