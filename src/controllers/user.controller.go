package controllers

import (
	"encoding/json"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	// Import the user package containing UserRepository and UserService
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	// Call the user service to get all users
	params := r.URL.Query()

	// set default values for page and per
	if !params.Has("page") { params.Set("page", "1") }
	if !params.Has("per") { params.Set("per", "10") }

	// fetch page and per from request query
	page, err_page := strconv.ParseInt(params.Get("page"), 10, 64)
	per, err_per := strconv.ParseInt(params.Get("per"), 10, 64)
	if err_page != nil || err_per != nil{
		http.Error(w, "Failed to parse request query params", http.StatusBadRequest)
		return
	}

	// get all users, no filters for now
	users, err := utills.GetUsers(db, bson.D{}, page - 1, per)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// CreateHandler handles the creation of a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
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

// SetDefaultBankAccountHandler handles the setting of a default bank account for a user
func SetDefaultBankAccountHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	type request struct {
		Username                string `json:"username"`
		DefaultBankAccountNumber string `json:"defaultBankAccountNumber"`
		DefaultBank             string `json:"defaultBank"`
	}
	// get user_id, default bank account number, default bank from request body
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	username := req.Username
	defaultBankAccountNumber := req.DefaultBankAccountNumber
	defaultBank := req.DefaultBank
	
	// Call the user service to set the default bank account
	err_str, err := utills.SetDefaultBankAccount(username, defaultBankAccountNumber, defaultBank, db)
	if err != nil {
		// show error message
		http.Error(w, err_str, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Default bank account set successfully")
}
