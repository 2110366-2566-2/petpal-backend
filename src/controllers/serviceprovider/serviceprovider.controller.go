package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	mail "petpal-backend/src/utills/email"
	svcp_utills "petpal-backend/src/utills/serviceprovider"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// GetSVCPsHandler handles the fetching of all service providers
func GetSVCPsHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
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

	// get all svcps, no filters for now
	svcps, err := svcp_utills.GetSVCPs(db, bson.D{}, page-1, per)
	if err != nil {
		http.Error(w, "Failed to get service providers", http.StatusInternalServerError)
		return
	}

	// remove sensitive data from svcps
	for i := range svcps {
		svcps[i].RemoveSensitiveData()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svcps)
}

// GetSVCPByIDHandler handles the fetching of a service provider by ID
func GetSVCPByIDHandler(w http.ResponseWriter, r *http.Request, db *models.MongoDB, id string) {
	svcp, err := svcp_utills.GetSVCPByID(db, id)
	if err != nil {
		http.Error(w, "Failed to get service provider", http.StatusInternalServerError)
		return
	}

	// remove sensitive data from svcps
	svcp.RemoveSensitiveData()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svcp)
}

func UpdateSVCPHandler(c *gin.Context, db *models.MongoDB, id string) {
	var svcp bson.M
	err := json.NewDecoder(c.Request.Body).Decode(&svcp)
	if err != nil {
		http.Error(c.Writer, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// check if the user is same as the svcp to be updated
	token, _ := c.Cookie("token")
	current_svcp, err := auth.GetCurrentSVCP(token, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}
	if id != current_svcp.SVCPID {
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = svcp_utills.UpdateSVCP(db, id, &svcp)
	if err != nil {
		http.Error(c.Writer, "Failed to update service provider", http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(svcp)
}

// RegisterHandler handles user registration
func RegisterSVCPHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	var createSVCP models.CreateSVCP
	if err := c.ShouldBindJSON(&createSVCP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Hash the password securely
	hashedPassword, err := auth.HashPassword(createSVCP.SVCPPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user instance
	createSVCP.SVCPPassword = hashedPassword
	newSVCP, err := auth.NewSVCP(createSVCP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create svcp "})
		return
	}

	// Insert the new user into the database
	newSVCP, err = svcp_utills.InsertSVCP(db, newSVCP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register svcp"})
		return
	}

	// Generate a JWT token
	tokenString, err := auth.GenerateToken(newSVCP.SVCPUsername, newSVCP.SVCPEmail, "svcp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token in cookies and send to frontend
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": tokenString})
}

// RegisterHandler handles user registration
func CurrentSVCPHandler(c *gin.Context, db *models.MongoDB) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	// Parse request body to get user data
	svcp, err := auth.GetCurrentSVCP(token, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get User Email request body :"+err.Error())
		return
	}

	// Remove sensitive data
	svcp.RemoveSensitiveData()
	
	// Set the content type header
	c.JSON(http.StatusAccepted, svcp)
}

func LoginSVCPHandler(c *gin.Context, db *models.MongoDB) {
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

func LogoutSVCPHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func ChangePassword(w http.ResponseWriter, r *http.Request, db *models.MongoDB) {
	type ChangePasswordReq struct {
		SVCPEmail   string `json:svcpemail`
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
	email := user.SVCPEmail
	newPassword := hashedPassword

	// Call the user service to set change password
	err_str, err := svcp_utills.ChangePassword(email, newPassword, db)
	if err != nil {
		// show error message
		http.Error(w, err_str, http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("set new password successfully")
}

func UploadDescriptionHandler(c *gin.Context, db *models.MongoDB) {
	var request models.SVCP
	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get current svcp
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	current_svcp, err := auth.GetCurrentSVCP(token, db)

	svcp_email := current_svcp.SVCPEmail
	description := request.Description

	err = svcp_utills.EditDescription(db, svcp_email, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Description uploaded successfully"})
}

func UploadSVCPLicenseHandler(c *gin.Context, db *models.MongoDB) {
	// Parse the form data, including the file upload
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	} 

	// Get the email from the current user
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	current_svcp, _ := auth.GetCurrentSVCP(token, db)
	email := current_svcp.SVCPEmail

	// Retrieve the uploaded file
	file, _, err := c.Request.FormFile("license")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Retrieving the File"})
		return
	}
	defer file.Close()

	// Read the file content as a byte slice
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file content"})
		return
	}
	err = svcp_utills.UploadSVCPLicense(db, fileContent, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send Confirmation email to the gmail
	emailSubject := "Petpal Service Provider Approved"
	emailContent := `
	<h4>สวัสดีครับ</h4>
	<p>ยินดีด้วย คุณยืนยันตัวตนกับทาง Petpal สำเร็จ!!!</p>
	`
	err = mail.SendEmailWithGmail(email, emailSubject, emailContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "update license successfull", "svcpEmail": email})
}

func AddServiceHandler(c *gin.Context, db *models.MongoDB) {
	var request struct {
		Service   models.Service `json:"service"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, login_err := c.Cookie("token")
	if login_err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+login_err.Error())
		return
	}
	current_svcp, _ := auth.GetCurrentSVCP(token, db)

	email := current_svcp.SVCPEmail
	err := svcp_utills.AddService(db, email, request.Service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service added successfully"})
}

func DeleteBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	// get current svcp
	token, login_err := c.Cookie("token")
	if login_err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+login_err.Error())
		return
	}
	current_svcp, _ := auth.GetCurrentSVCP(token, db)

	err := svcp_utills.DeleteBankAccount(db, current_svcp.SVCPEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank account deleted successfully"})
}

func SetDefaultBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	var request struct {
		DefaultAccountNumber string `json:"defaultAccountNumber"`
		DefaultBank          string `json:"defaultBank"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get current svcp
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return
	}
	current_svcp, err := auth.GetCurrentSVCP(token, db)

	svcp_email := current_svcp.SVCPEmail
	default_bank_account := request.DefaultAccountNumber
	default_bank := request.DefaultBank

	_, err = svcp_utills.SetDefaultBankAccount(svcp_email, default_bank_account, default_bank, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default bank account set successfully"})
}
