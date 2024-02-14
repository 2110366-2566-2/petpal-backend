package controllers

import (
	"encoding/json"
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

// GetUsersHandler godoc
//
// @Summary		Get all users
// @Description	Get all users (authentication not required)
// @Tags		user
//
// @Accept		json
// @Produce		json
// @Param		page	query		int	false	"page"
// @Param		per		query		int	false	"per"
//
// @Success		200		{object}	[]models.User
// @Failure		400		{object}	string
// @Failure		500		{object}	string
//
// @Router		/users [get]
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
	var user bson.M
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the user service to update the user
	err_str, err := utills.UpdateUser(db, &user, c.Param("id"))
	if err != nil {
		// show error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// RegisterUserHandler godoc
// @Summary Register User
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param requestBody body models.CreateUser
// @Success 200 {object} RegisterUserResp
// @Failure 400 {object} string
// @Failure 500 {object} string
//
//	@Router			/user/register [post]
type RegisterUserResp struct {
	// Define the 10 fields here
	massage string
	token   string
}

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
	newUser, err := auth.NewUser(createUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user "})
		return
	}

	// Insert the new user into the database
	newUser, err = user_utills.InsertUser(db, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user" + err.Error()})
		return
	}

	// Generate a JWT token
	tokenString, err := auth.GenerateToken(newUser.Username, newUser.Password, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token in cookies and send to frontend
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": tokenString})
}

// CurrentUserHandler godoc
// @Summary Get Current User
// @Description Get the details of the currently authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 202 {object} models.User "User details"
// @Failure 400 {string} string "Failed to get token from Cookie plase login first"
// @Failure 500 {string} string "Failed to get User Email request body"
// @Router /user/me [get]
func CurrentUserHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	user, err := auth.GetCurrentUserByGinContext(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	// Set the content type header
	c.JSON(http.StatusAccepted, user)
}

func LoginUserHandler(c *gin.Context, db *models.MongoDB) {
	var user models.LoginReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := auth.Login(db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", u.AccessToken, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, u)
}

func LogoutUserHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func GetUserPetsHandler(c *gin.Context, db *models.MongoDB) {
	pets, err := user_utills.GetUserPet(db, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pets": pets})
}

func AddUserPetHandler(c *gin.Context, db *models.MongoDB) {
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err_str, err := user_utills.AddUserPet(db, &pet, c.Param("id"))
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
func UpdateUserPetHandler(c *gin.Context, db *models.MongoDB) {
	pet_idx, err := strconv.Atoi(c.Param("idx"))
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

	err_str, err := user_utills.UpdateUserPet(db, &pet, c.Param("id"), pet_idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet updated successfully"})
}

func DeleteUserPetHandler(c *gin.Context, db *models.MongoDB) {
	pet_idx, err := strconv.Atoi(c.Param("idx"))
	if err != nil || pet_idx < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse pet index"})
		return
	}
	err_str, err := user_utills.DeleteUserPet(db, c.Param("id"), pet_idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err_str})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pet deleted successfully"})

}

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

func ChangePassword(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	type ChangePasswordReq struct {
		UserEmail   string `json:useremail`
		NewPassword string `json:newpassword`
	}
	var user ChangePasswordReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(user.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusBadRequest)
		return
	}
	email := user.UserEmail
	newPassword := hashedPassword

	// Call the user service to set change password
	err_str, err := utills.ChangePassword(email, newPassword, db)
	if err != nil {
		// show error message
		http.Error(w, err_str, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("set new password successfully")
}

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
