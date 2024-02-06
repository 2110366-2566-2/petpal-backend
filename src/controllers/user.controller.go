package controllers

import (
	"encoding/json"
	"net/http"
	"petpal-backend/models"
	// Import the user package containing UserRepository and UserService
)

type UserController struct {
	db models.MongoDB
}

// CreateHandler handles the creation of a new user
func (uc *UserController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get user data
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Call the user service to create a new user
	createdUser, err := uc.UserService.RegisterUser(newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}
