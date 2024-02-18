package controllers

import (
	"encoding/json"
	"errors"
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

func UpdateSVCPHandler(c *gin.Context, db *models.MongoDB) {
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}
	var svcp bson.M
	if err := c.ShouldBindJSON(&svcp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = svcp_utills.UpdateSVCP(db, current_svcp.SVCPID, &svcp)
	if err != nil {
		http.Error(c.Writer, "Failed to update service provider", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "SVCP updated successfully"})
}

func UploadDescriptionHandler(c *gin.Context, db *models.MongoDB) {
	var request models.SVCP
	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
		Service models.Service `json:"service"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := current_svcp.SVCPEmail
	err = svcp_utills.AddService(db, email, request.Service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service added successfully"})
}

func DeleteBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = svcp_utills.DeleteBankAccount(db, current_svcp.SVCPEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank account deleted successfully"})
}

func SetDefaultBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request struct {
		DefaultAccountNumber string `json:"defaultAccountNumber"`
		DefaultBank          string `json:"defaultBank"`
	}

	err = json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func _authenticate(c *gin.Context, db *models.MongoDB) (*models.SVCP, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get token from Cookie plase login first, "+err.Error())
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.User:
		err = errors.New("need token of type SVCP but recives token SVCP type")
		c.JSON(http.StatusBadRequest, err.Error())
		return nil, nil
		// Handle user
	case *models.SVCP:
		return entity, nil
		// Handle svcp
	}
	err = errors.New("need token of type SVCP but wrong type")
	c.JSON(http.StatusBadRequest, err.Error())
	return nil, err
}
