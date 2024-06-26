package controllers

import (
	"errors"
	"io"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	"petpal-backend/src/utills/chat/chathistory"
	user_utills "petpal-backend/src/utills/user"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	// Import the user package containing UserRepository and UserService
)

// GetUsersHandler godoc
//
// @Summary     Get all user
// @Description Retrieve all user
// @Tags        User
//
// @Accept      json
// @Produce     json
//
// @Param       page      query    int    false        "Page number"
// @Param       per       query    int    false        "Number of users per page"
//
// @Success     200      {array} 	models.User    				"Success"
// @Failure     400      {object} 	models.BasicErrorRes      	"Bad request"
// @Failure     500      {object} 	models.BasicErrorRes		"Internal server error"
//
// @Router      /user [get]
func GetUsersHandler(c *gin.Context, db *models.MongoDB) {
	// Call the user service to get all users
	params := c.Request.URL.Query()

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
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to parse page and per"})
		return
	}

	// get all users, no filters for now
	users, err := user_utills.GetUsers(db, bson.D{}, page-1, per)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByIDHandler godoc
//
// @Summary     Get user by ID
// @Description Retrieve user information by ID
// @Tags        User
//
// @Accept      json
// @Produce     json
//
// @Param       id      path    string    true        "User ID"
//
// @Success     200      {object} models.User    			"Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/{id} [get]
func GetUserByIDHandler(c *gin.Context, db *models.MongoDB, id string) {
	// Call the user service to get a user by email
	user, err := user_utills.GetUserByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserHandler godoc
//
// @Summary     Update user information
// @Description Update user information (authentication required)
// @Tags        User
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       user      body    models.User       true        "User object that needs to be updated"
//
// @Success     200      {object} models.BasicRes    "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user [put]
func UpdateUserHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	var userUpdate bson.M
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	// Call the user service to update the user
	err_str, err := user_utills.UpdateUser(db, &userUpdate, currentUser.ID)
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err_str})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, models.BasicRes{Message: "User updated successfully"})
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
// @Success     200      {object} object{pets=[]models.Pet}    "Success"
// @Failure     400      {object} models.BasicErrorRes  "Bad request"
// @Failure     401      {object} models.BasicErrorRes  "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes  "Internal server error"
//
// @Router      /user/pets/{id} [get]
func GetUserPetsByIdHandler(c *gin.Context, db *models.MongoDB, id string) {
	pets, err := user_utills.GetUserPet(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
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
// @Success     200      {object} object{pets=[]models.Pet}    "Success"
// @Failure     401      {object} models.BasicErrorRes  "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes  "Internal server error"
//
// @Router      /user/pets/me [get]
func GetCurrentUserPetsHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}
	pets, err := user_utills.GetUserPet(db, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
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
// @Param       pet      body    models.CreatePet       true        "Pet object to be added"
//
// @Success     200      {object} models.BasicRes    "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/pets [post]
func AddUserPetHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	var pet models.CreatePet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	err_str, err := user_utills.AddUserPet(db, &pet, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err_str})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Pet added successfully"})
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
// @Success     200      {object} models.BasicRes		    "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     404      {object} models.BasicErrorRes      "Not found"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/pets/{idx} [put]
func UpdateUserPetHandler(c *gin.Context, db *models.MongoDB, idx string) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	pet_idx, err := strconv.Atoi(idx)
	if err != nil || pet_idx < 0 {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to parse pet index"})
		return
	}

	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	err_str, err := user_utills.UpdateUserPet(db, &pet, currentUser.ID, pet_idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err_str})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Pet updated successfully"})
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
// @Success     200      {object} models.BasicRes		   "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     404      {object} models.BasicErrorRes      "Not found"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/pets/{idx} [delete]
func DeleteUserPetHandler(c *gin.Context, db *models.MongoDB, idx string) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	pet_idx, err := strconv.Atoi(idx)
	if err != nil || pet_idx < 0 {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to parse pet index"})
		return
	}
	err_str, err := user_utills.DeleteUserPet(db, currentUser.ID, pet_idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err_str})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Pet deleted successfully"})
}

// SetDefaultBankAccountHandler godoc
//
// @Summary     Set default bank account
// @Description Set default bank account for the current user
// @Tags        User
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       default_account_number      body    string    true    "Default account number"
// @Param       default_bank                body    string    true    "Default bank"
//
// @Success     200      {object} models.BasicRes		    "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/set-default-bank-account [post]
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
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	// Call the user service to set the default bank account
	err_str, err := user_utills.SetDefaultBankAccount(currentUser.Email, req.DefaultAccountNumber, req.DefaultBank, db)
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err_str})
		return
	}
	c.JSON(http.StatusOK, models.BasicRes{Message: "Default bank account set successfully"})
}

// DeleteBankAccountHandler godoc
//
// @Summary     Delete bank account
// @Description Delete the bank account associated with the current user
// @Tags        User
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Success     200      {object} models.BasicRes    "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/delete-bank-account [delete]
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
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err_str})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, models.BasicRes{Message: "Bank account deleted successfully"})
}

// UploadImageHandler godoc
//
// @Summary     Upload profile image
// @Description Uploads a profile image for the current user
// @Tags        User
//
// @Accept      multipart/form-data
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       profileImage      formData    file      true        "Profile image file to upload"
//
// @Success     202      {object} models.BasicRes    "Accepted"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/uploadProfileImage [post]
func UploadImageHandler(c *gin.Context, db *models.MongoDB) {
	// Parse the multipart form data
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		// If unable to parse the form, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Error Parsing the Form"})
		return
	}

	// Retrieve the uploaded file
	file, _, err := c.Request.FormFile("profileImage")
	if err != nil {
		// If an error occurs while retrieving the file, respond with a bad request and error message
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Error Retrieving the File"})
		return
	}
	defer file.Close()

	// Read the content of the uploaded file
	fileContent, err := io.ReadAll(file)
	if err != nil {
		// If there is an error reading the file content, respond with a internal server error and error message
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Error Reading the File"})
		return
	}

	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
	}
	switch entity := entity.(type) {
	case *models.User:
		// Perform the upload of the profile image to the database using a utility function
		response, err := user_utills.UploadProfileImage(entity.Email, fileContent, "user", db)
		if err != nil {
			// If there is an error during the profile image upload, respond with an internal server error and error message
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		c.JSON(http.StatusAccepted, response)
		// Handle user
	case *models.SVCP:
		response, err := user_utills.UploadProfileImage(entity.SVCPEmail, fileContent, "svcp", db)
		if err != nil {
			// If there is an error during the profile image upload, respond with an internal server error and error message
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		c.JSON(http.StatusAccepted, response)
	}
}

// GetSearchHistoryHandler godoc
//
// @Summary     Get search history
// @Description Gets the search history of the current user
// @Tags        User
//
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Success     200      {object} models.UserSearchHistory  "Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/get-search-history [get]
func GetSearchHistoryHandler(c *gin.Context, db *models.MongoDB) {
	currentUser, err := _authenticate(c, db)
	if err != nil {
		return
	}

	search_history, err := user_utills.GetSearchHistory(db, currentUser.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserSearchHistory{User: *currentUser, SearchHistory: search_history})
}

// GetChatsHandler godoc
//
// @Summary 	Get the current user's chats
// @Description Get chat rooms of the current user. Each `Chat.messages` contains only *one* latest message.
// @Tags 		User
//
// @Security ApiKeyAuth
//
// @Produce  	json
//
// @Param 		page	query	int 	false	"Page number of chat rooms (default 1)"
// @Param 		per 	query	int 	false 	"Number of chat rooms per page (default 10)"
//
// @Success 	200      {object} []models.Chat         "Success"
// @Failure     400      {object} models.BasicErrorRes  "Bad request"
// @Failure     401      {object} models.BasicErrorRes  "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes  "Internal server error"
//
// @Router /user/chats [get]
func GetChatsHandler(c *gin.Context, db *models.MongoDB) {
	current_user, err := _authenticate(c, db)
	if err != nil {
		return
	}

	params := c.Request.URL.Query()

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
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid page or per number"})
		return
	}

	chatHistory, err := chathistory.GetChatsById(db, current_user.ID, page, per, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chatHistory)
}

func _authenticate(c *gin.Context, db *models.MongoDB) (*models.User, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.User:
		return entity, nil
		// Handle user
	case *models.SVCP:
		err = errors.New("need token of type User but recives token SVCP type")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, err
		// Handle svcp
	}
	err = errors.New("need token of type User but wrong type")
	c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
	return nil, err
}
