package controllers

import (
	"encoding/json"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"

	"go.mongodb.org/mongo-driver/mongo"
	// Import the user package containing UserRepository and UserService
)

// CreateHandler handles the creation of a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request, db *mongo.Client) {
	// Parse request body to get user data
	var createNewUser models.CreateUser
	err := json.NewDecoder(r.Body).Decode(&createNewUser)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Call the user service to create a new user
	createdUser, err := utills.NewUser(createNewUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}
