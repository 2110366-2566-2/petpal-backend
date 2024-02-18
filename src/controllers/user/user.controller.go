package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	user_utills "petpal-backend/src/utills/user"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	// Import the user package containing UserRepository and UserService
)

// GetUsersHandler godoc
//
// @Summary     Get all users
// @Description Retrieve all users
// @Tags        Users
//
// @Accept      json
// @Produce     json
//
// @Param       page      query    int    false        "Page number"
// @Param       per       query    int    false        "Number of users per page"
//
// @Success     200      {array} User    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /user [get]
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
	users, err := user_utills.GetUsers(db, bson.D{}, page-1, per)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByIDHandler godoc
//
// @Summary     Get user by ID
// @Description Retrieve user information by ID
// @Tags        Users
//
// @Accept      json
// @Produce     json
//
// @Param       id      path    string    true        "User ID"
//
// @Success     200      {object} User    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /user/{id} [get]
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB, id string) {
	// Call the user service to get a user by email
	user, err := user_utills.GetUserByID(db, id)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler godoc
//
// @Summary     Update user information
// @Description Update user information (authentication required)
// @Tags        Users
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       user      body    models.User       true        "User object that needs to be updated"
//
// @Success     200      {object} object{message=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     401      {object} object{error=string}      "Unauthorized"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /user [put]
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
	err_str, err := user_utills.UpdateUser(db, &user, currentUser.ID)
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// GetUserPetsByIdHandler godoc
//
// @Summary     Get user's pets by user ID
// @Description Get user's pets by user ID
// @Tags        Pets
//
// @Accept      json
// @Produce     json
//
// @Param       id      path    string    true        "User ID"
//
// @Success     200      {object} object{pets=[]Pet}    "Success"
// @Failure     400      {object} object{error=string}  "Bad request"
// @Failure     401      {object} object{error=string}  "Unauthorized"
// @Failure     500      {object} object{error=string}  "Internal server error"
//
// @Router      /users/pets/{id} [get]
func GetUserPetsByIdHandler(c *gin.Context, db *models.MongoDB, id string) {
	pets, err := user_utills.GetUserPet(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pets": pets})
}

// GetCurrentUserPetsHandler godoc
//
// @Summary     Get current user's pets
// @Description Get current user's pets (authentication required)
// @Tags        Pets
//
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Success     200      {object} object{pets=[]Pet}    "Success"
// @Failure     401      {object} object{error=string}  "Unauthorized"
// @Failure     500      {object} object{error=string}  "Internal server error"
//
// @Router      /user/pets/me [get]
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

// AddUserPetHandler godoc
//
// @Summary     Add a new pet for the user
// @Description Add a new pet for the user (authentication required)
// @Tags        Pets
// @id		  	AddUserPetHandler
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       pet      body    models.Pet       true        "Pet object to be added"
//
// @Success     200      {object} object{message=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     401      {object} object{error=string}      "Unauthorized"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      user/pets [post]
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

// UpdateUserPetHandler godoc
//
// @Summary     Update user's pet information
// @Description Update user's pet information (authentication required)
// @Tags        Pets
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       idx      path    string    true        "Pet Index"
// @Param       pet      body    models.Pet       true        "Pet object that needs to be updated"
//
// @Success     200      {object} object{message=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     401      {object} object{error=string}      "Unauthorized"
// @Failure     404      {object} object{error=string}      "Not found"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      user/pets/{idx} [put]
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

// DeleteUserPetHandler godoc
//
// @Summary     Delete user's pet
// @Description Delete user's pet (authentication required)
// @Tags        Pets
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       idx      path    string    true        "Pet Index"
//
// @Success     200      {object} object{message=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     401      {object} object{error=string}      "Unauthorized"
// @Failure     404      {object} object{error=string}      "Not found"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /user/pets/{idx} [delete]
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
func SetDefaultBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}
	type SetDefaultBankAccountReq struct {
		DefaultAccountNumber string
		DefaultBank          string
	}
	var req SetDefaultBankAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the user service to set the default bank account
	err_str, err := user_utills.SetDefaultBankAccount(currentUser.Email, req.DefaultAccountNumber, req.DefaultBank, db)
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Default bank account set successfully"})
}

// DeleteBankAccountHandler handles the deletion of a bank account for a user
func DeleteBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	// get user_id from request body
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	// Call the user service to delete the bank account
	err_str, err := user_utills.DeleteBankAccount(currentUser.Email, db)
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Bank account deleted successfully"})
}

// UploadImageHandler handles the HTTP request for uploading a profile image.
func UploadImageHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}
	// Parse the multipart form data
	err = c.Request.ParseMultipartForm(10 << 20)
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

	// Read the content of the uploaded file
	fileContent, err := io.ReadAll(file)
	if err != nil {
		// If there is an error reading the file content, respond with a internal server error and error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file content"})
		return
	}

	// Perform the upload of the profile image to the database using a utility function
	response, err := user_utills.UploadProfileImage(currentUser.Email, fileContent, "user", db)
	if err != nil {
		// If there is an error during the profile image upload, respond with an internal server error and error message
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// If everything is successful, respond with an accepted status and the response
	c.JSON(http.StatusAccepted, response)
}

// GetProfileImageHandler godoc
//
// @Summary     Get profile image for a user
// @Description Get profile image for a user
// @Tags        Users
//
// @Accept      json
// @Produce     json
//
// @Param       email    query    string    true    "User email"
// @Param       userType path      string    true    "User type (e.g., 'customer' or 'provider')"
//
// @Success     202      {object} object{image=string}    "Accepted"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /user/{userType}/profile-image [get]
func GetProfileImageHandler(c *gin.Context, userType string, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	// // Perform the upload of the profile image to the database using a utility function
	response, err := user_utills.GetProfileImage(currentUser.Email, userType, db)
	if err != nil {
		// If there is an error during the profile image upload, respond with an internal server error and error message
		c.JSON(http.StatusInternalServerError, err.Error())
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
		err = errors.New("need token of type User but recives token SVCP type")
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, nil
		// Handle svcp
	}
	err = errors.New("need token of type User but wrong type")
	c.JSON(http.StatusBadRequest, err.Error())
	return nil, err
}
