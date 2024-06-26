package utills

import (
	"context"
	"fmt"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserPet(db *models.MongoDB, id string) (*[]models.Pet, error) {
	user, err := GetUserByID(db, id)
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

func CreateNewPet(pet *models.CreatePet, ownerUsername string) *models.Pet {
	return &models.Pet{
		OwnerUsername:     ownerUsername,
		Name:              pet.Name,
		Gender:            pet.Gender,
		Age:               pet.Age,
		Pet_type:          pet.Pet_type,
		HealthInformation: pet.HealthInformation,
		Certificate:       pet.Certificate,
		BehaviouralNotes:  pet.BehaviouralNotes,
		Vaccinations:      pet.Vaccinations,
		DietyPreferences:  pet.DietyPreferences,
		Breed:             pet.Breed,
	}
}

func AddUserPet(db *models.MongoDB, createPet *models.CreatePet, user *models.User) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by email
	user_id := user.ID
	pet := CreateNewPet(createPet, user.Username)

	filter := bson.D{{Key: "_id", Value: user_id}}
	res, err := user_collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$push", Value: bson.D{{Key: "pets", Value: pet}}}})
	if res == nil {
		return "User not found (id=" + user_id + ")", err
	}
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
	filter := bson.D{{Key: "_id", Value: user_id}}
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
	filter := bson.D{{Key: "_id", Value: user_id}}
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
