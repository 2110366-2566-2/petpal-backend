package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"

	"github.com/gin-gonic/gin"
	// Import the user package containing UserRepository and UserService
)

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

// UploadImageHandler handles the HTTP request for uploading a profile image.
func UploadImageHandler(c *gin.Context, db *models.MongoDB) {
	// Parse the multipart form data
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		// If unable to parse the form, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}

	// Retrieve the uploaded file
	file, _, err := c.Request.FormFile("profileImage")
	if err != nil {
		// If an error occurs while retrieving the file, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Retrieving the File"})
		return
	}
	defer file.Close()

	// Retrieve the username from the form data
	username := c.Request.FormValue("username")

	// Check if the username is empty
	if username == "" {
		// If username is empty, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Read the content of the uploaded file
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		// If there is an error reading the file content, respond with a internal server error and error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file content"})
		return
	}

	// Perform the upload of the profile image to the database using a utility function
	response, err := utills.UploadProfileImage(username, fileContent, db)
	if err != nil {
		// If there is an error during the profile image upload, respond with an internal server error and error message
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// If everything is successful, respond with an accepted status and the response
	c.JSON(http.StatusAccepted, response)
}
