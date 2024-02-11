package controllers

import (
	"encoding/json"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"
	"petpal-backend/src/utills/auth"

	"github.com/gin-gonic/gin"
	// Import the user package containing UserRepository and UserService
)

// RegisterHandler handles user registration
func RegisterUserHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	var createUser models.CreateUser
	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password securely
	hashedPassword, err := auth.HashPassword(createUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user instance
	createUser.Password = hashedPassword
	newUser, err := utills.NewUser(createUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user "})
		return
	}

	// Insert the new user into the database
	newUser, err = utills.InsertUser(db, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Generate a JWT token
	tokenString, err := auth.GenerateToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token in cookies and send to frontend
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": tokenString})
}

// SetDefaultBankAccountHandler handles the setting of a default bank account for a user
func SetDefaultBankAccountHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	type request struct {
		Username                 string `json:"username"`
		DefaultBankAccountNumber string `json:"defaultBankAccountNumber"`
		DefaultBank              string `json:"defaultBank"`
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
