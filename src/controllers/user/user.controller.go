package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	user_utills "petpal-backend/src/utills/user"
	utills "petpal-backend/src/utills/user"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	// Import the user package containing UserRepository and UserService
)

// GetUsersHandler handles the fetching of all users
func GetUsersHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	// Call the user service to get all users
	params := r.URL.Query()

	// set default values for page and per
	if !params.Has("page") {
		params.Set("page", "1")
	}
	if !params.Has("per") {
		params.Set("per", "10")
	}

	// fetch page and per from request query
	page, err_page := strconv.ParseInt(params.Get("page"), 10, 64)
	per, err_per := strconv.ParseInt(params.Get("per"), 10, 64)
	if err_page != nil || err_per != nil {
		http.Error(w, "Failed to parse request query params", http.StatusBadRequest)
		return
	}

	// get all users, no filters for now
	users, err := utills.GetUsers(db, bson.D{}, page-1, per)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByIDHandler handles the fetching of a user by id
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB, id string) {
	// Call the user service to get a user by email
	user, err := utills.GetUserByID(db, id)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	var user bson.M
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the user service to update the user
	err_str, err := utills.UpdateUser(db, &user, currentUser.ID)
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

}

// GetUserPetsHandler for get list of user's pet
func GetUserPetsByIdHandler(c *gin.Context, db *models.MongoDB, id string) {
	pets, err := user_utills.GetUserPet(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pets": pets})
}

// GetUserPetsHandler for get list of user's pet
func GetCurrentUserPetsHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}
	pets, err := user_utills.GetUserPet(db, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pets": pets})
}

func AddUserPetHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err_str, err := user_utills.AddUserPet(db, &pet, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet added successfully"})
}

// UpdateUserPetHandler for updating user's pet
//
// note: the body of the request should contain all of the updated pet's details
// otherwise the missing fields will be set to their zero values
// also: this updates the pet at the specified index param `idx`
func UpdateUserPetHandler(c *gin.Context, db *models.MongoDB, idx string) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	pet_idx, err := strconv.Atoi(idx)
	if err != nil || pet_idx < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse pet index"})
		return
	}

	// this binding sets missing fields to their zero values
	// the pet model does not have any validation tags
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err_str, err := user_utills.UpdateUserPet(db, &pet, currentUser.ID, pet_idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet updated successfully"})
}

func DeleteUserPetHandler(c *gin.Context, db *models.MongoDB, idx string) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	pet_idx, err := strconv.Atoi(idx)
	if err != nil || pet_idx < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse pet index"})
		return
	}
	err_str, err := user_utills.DeleteUserPet(db, currentUser.ID, pet_idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet deleted successfully"})

}

// SetDefaultBankAccountHandler handles the setting of a default bank account for a user
func SetDefaultBankAccountHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	// get user_id, default bank account number, default bank from request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	email := user.Email
	defaultAccountNumber := user.DefaultAccountNumber
	defaultBank := user.DefaultBank

	// Call the user service to set the default bank account
	err_str, err := utills.SetDefaultBankAccount(email, defaultAccountNumber, defaultBank, db)
	if err != nil {
		// show error message
		http.Error(w, err_str, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Default bank account set successfully")
}

// DeleteBankAccountHandler handles the deletion of a bank account for a user
func DeleteBankAccountHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	// get user_id from request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	email := user.Email

	// Call the user service to delete the bank account
	err_str, err := utills.DeleteBankAccount(email, db)
	if err != nil {
		// show error message
		http.Error(w, err_str, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Bank account deleted successfully")
}

// UploadImageHandler handles the HTTP request for uploading a profile image.
func UploadImageHandler(c *gin.Context, userType string, db *models.MongoDB) {
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

	// Retrieve the email from the form data
	email := c.Request.FormValue("email")

	// Check if the email is empty
	if email == "" {
		// If email is empty, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
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
	response, err := utills.UploadProfileImage(email, fileContent, userType, db)
	if err != nil {
		// If there is an error during the profile image upload, respond with an internal server error and error message
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// If everything is successful, respond with an accepted status and the response
	c.JSON(http.StatusAccepted, response)
}

// UploadImageHandler handles the HTTP request for uploading a profile image.
func GetProfileImageHandler(c *gin.Context, userType string, db *models.MongoDB) {

	// Retrieve the email from the form data
	email := c.Request.FormValue("email")

	// Check if the email is empty
	if email == "" {
		// If email is empty, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// // Perform the upload of the profile image to the database using a utility function
	response, err := utills.GetProfileImage(email, userType, db)
	if err != nil {
		// If there is an error during the profile image upload, respond with an internal server error and error message
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// If everything is successful, respond with an accepted status and the response
	c.JSON(http.StatusAccepted, response)
}

func _authenticate(c *gin.Context, db *models.MongoDB) (*models.User, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.User:
		return entity, nil
		// Handle user
	case *models.SVCP:
		err = errors.New("Need token of type User but recives token SVCP type")
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, nil
		// Handle svcp
	}
	err = errors.New("Need token of type User but wrong type")
	c.JSON(http.StatusBadRequest, err.Error())
	return nil, err
}
