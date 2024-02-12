package controllers

import (
	"encoding/json"
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
	newUser, err := auth.NewUser(createUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user "})
		return
	}

	// Insert the new user into the database
	newUser, err = user_utills.InsertUser(db, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
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

// RegisterHandler handles user registration
func CurrentUserHandler(c *gin.Context, db *models.MongoDB) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	// Parse request body to get user data
	user, err := auth.GetCurrentUser(token, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get User Email request body :"+err.Error())
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

// LogoutUserHandler
func LogoutUserHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// GetUserPetsHandler for get list of user's pet
func GetUserPetsHandler(c *gin.Context, db *models.MongoDB) {
	type GetUserPetReq struct {
		UserEmail string `json:useremail`
	}
	var user GetUserPetReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pets, err := user_utills.GetUserPet(db, user.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pets": pets, "useremail": user.UserEmail})
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