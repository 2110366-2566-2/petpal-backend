package controllers

import (
	"encoding/json"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	utills "petpal-backend/src/utills/user"

	"github.com/gin-gonic/gin"
	// Import the user package containing UserRepository and UserService
)

// ChangePassword godoc
//
// @Summary     Change user password
// @Description Change user password
// @Tags        Users
//
// @Accept      json
// @Produce     json
//
// @Param       user_email     body    string    true        "User email"
// @Param       new_password   body    string    true        "New password"
// @Param       login_type     body    string    true        "Login type"
//
// @Success     200      {object} object{message=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /change-password [post]
func ChangePassword(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	type ChangePasswordReq struct {
		UserEmail   string
		NewPassword string
		LoginType   string
	}
	var changePasswordReq ChangePasswordReq
	err := json.NewDecoder(r.Body).Decode(&changePasswordReq)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(changePasswordReq.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusBadRequest)
		return
	}
	email := changePasswordReq.UserEmail

	// Call the user service to set change password
	err_str, err := utills.ChangePassword(email, hashedPassword, db)
	if err != nil {
		// show error message
		http.Error(w, err_str, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("set new password successfully")
}

// GetCurrentEntityHandler godoc
//
// @Summary     Get current entity
// @Description Get the current entity based on the provided context
// @Tags        Authentication
//
// @Accept      json
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Success     200      {object} models.User     			"Accepted"
// @Failure     400      {object} object{error=string}      "Bad request"
//
// @Router      /current-entity [get]
func GetCurrentEntityHandler(c *gin.Context, db *models.MongoDB) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	c.JSON(http.StatusAccepted, entity)
}

// LoginHandler godoc
//
// @Summary     User login
// @Description Authenticate user and generate access token
// @Tags        Authentication
//
// @Accept      json
// @Produce     json
//
// @Param       loginReq        body    models.LoginReq    true    "login request"
//
// @Success     200      {object} object    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /login [post]
func LoginHandler(c *gin.Context, db *models.MongoDB) {
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

// LogoutHandler godoc
//
// @Summary     User logout
// @Description Invalidate user session and clear access token
// @Tags        Authentication
//
// @Produce     json
//
// @Success     200      {object} object{message=string}    "Success"
//
// @Router      /logout [post]
func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// RegisterUserHandler godoc
//
// @Summary     User registration
// @Description Register a new user
// @Tags        Authentication
//
// @Accept      json
// @Produce     json
//
// @Param       user_data     body    models.CreateUser    true    "User registration data"
//
// @Success     200      {object} object{message=string, token=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /register-user [post]
func RegisterUserHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	var createUser models.CreateUser

	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Generate a JWT token
	tokenString, err := auth.RegisterUser(createUser, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set token in cookies and send to frontend
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": tokenString})
}

// RegisterSVCPHandler godoc
//
// @Summary     Service provider registration
// @Description Register a new service provider
// @Tags        Authentication
//
// @Accept      json
// @Produce     json
//
// @Param       svcp_data     body    models.CreateSVCP    true    "Service provider registration data"
//
// @Success     200      {object} object{message=string, token=string}    "Success"
// @Failure     400      {object} object{error=string}      "Bad request"
// @Failure     500      {object} object{error=string}      "Internal server error"
//
// @Router      /register-svcp [post]
func RegisterSVCPHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	var createUser models.CreateSVCP
	if err := c.ShouldBindJSON(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Generate a JWT token
	tokenString, err := auth.RegisterSVCP(createUser, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set token in cookies and send to frontend
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": tokenString})
}
