package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	"petpal-backend/src/utills/chat/chathistory"
	mail "petpal-backend/src/utills/email"
	svcp_utills "petpal-backend/src/utills/serviceprovider"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

// GetSVCPsHandler godoc
//
// @Summary 	Get all service providers
// @Description Get all service providers (authentication not required) and sensitive information is censorred
// @Tags 		ServiceProviders
//
// @Accept  	json
// @Produce  	json
//
// @Param 		page	query	int 	false	"Page number(default 1)"
// @Param 		per 	query	int 	false 	"Number of items per page(default 10)"
//
// @Success 200 {array} models.SVCP
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /serviceproviders [get]
func GetSVCPsHandler(c *gin.Context, db *models.MongoDB) {
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
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid page or per"})
		return
	}

	// get all svcps, no filters for now
	svcps, err := svcp_utills.GetSVCPs(db, bson.D{}, page-1, per)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to get service providers" + err.Error()})
		return
	}

	// remove sensitive data from svcps
	for i := range svcps {
		svcps[i].RemoveSensitiveData()
	}

	c.JSON(http.StatusOK, svcps)
}

// GetSVCPByIDHandler godoc
//
// @Summary 	Get service provider by ID
// @Description Get service provider by ID (authentication not required) and sensitive information is censorred
// @Tags 		ServiceProviders
//
// @Accept  	json
// @Produce  	json
//
// @Param 		id	path	string 	true	"Service Provider ID"
//
// @Success 200 {object} models.SVCP
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /serviceproviders/{id} [get]
func GetSVCPByIDHandler(c *gin.Context, db *models.MongoDB, id string) {
	svcp, err := svcp_utills.GetSVCPByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to get service provider"})
		return
	}

	// remove sensitive data from svcps
	svcp.RemoveSensitiveData()

	c.JSON(http.StatusOK, svcp)
}

// UpdateSVCPHandler godoc
//
// @Summary 	Update service provider
// @Description Update service provider (authentication required and only the service provider can update their own profile)
// @Tags 		ServiceProviders
//
// @Security ApiKeyAuth
//
// @Accept  	json
// @Produce  	json
//
// @Param 		svcp	body 	object 	true	"Service Provider Object (only the fields to be updated)"
//
// @Success 200 {object} models.BasicRes
// @Failure 400 {object} models.BasicErrorRes
//
// @Router /serviceproviders [put]
func UpdateSVCPHandler(c *gin.Context, db *models.MongoDB) {
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}
	var svcp bson.M
	if err := c.ShouldBindJSON(&svcp); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid request"})
		return
	}

	err = svcp_utills.UpdateSVCP(db, current_svcp, &svcp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to update service provider"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, models.BasicRes{Message: "Service provider updated successfully"})
}

// UploadDescriptionHandler godoc
//
// @Summary 	Upload service provider description
// @Description Upload service provider description (authentication required and only the service provider can update their own profile)
// @Tags 		ServiceProviders
//
// @Security ApiKeyAuth
//
// @Accept  	json
// @Produce  	json
//
// @Param 		object	body 	object{description=string} 	true	"Description Request Object"
//
// @Success 200 {object} models.BasicRes
// @Failure 400 {object} models.BasicErrorRes
//
// @Router /serviceproviders/upload-description [post]
func UploadDescriptionHandler(c *gin.Context, db *models.MongoDB) {
	var request models.SVCP
	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid request"})
		return
	}

	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	svcp_email := current_svcp.SVCPEmail
	description := request.Description

	err = svcp_utills.EditDescription(db, svcp_email, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Description updated successfully"})
}

// UploadSVCPLicenseHandler godoc
//
// @Summary 	Upload service provider license
// @Description Upload service provider license (authentication required)
// @Tags 		ServiceProviders
//
// @Security ApiKeyAuth
//
// @Accept  	multipart/form-data
// @Produce  	json
//
// @Param 		license	formData 	file 	true	"License File"
//
// @Success 200 {object} object{message=string,svcpEmail=string}
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /serviceproviders/upload-license [post]
func UploadSVCPLicenseHandler(c *gin.Context, db *models.MongoDB) {
	// Parse the form data, including the file upload
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Error Parsing the Form"})
		return
	}

	// Get the email from the current user
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	email := current_svcp.SVCPEmail

	// Retrieve the uploaded file
	file, _, err := c.Request.FormFile("license")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Error Retrieving the File"})
		return
	}
	defer file.Close()

	// Read the file content as a byte slice
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Error Reading the File"})
		return
	}
	err = svcp_utills.UploadSVCPLicense(db, fileContent, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
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
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "update license successfull", "svcpEmail": email})
}

// AddServiceHandler godoc
//
// @Summary 	Add service
// @Description Add service (authentication required and only the service provider can update their own profile)
// @Tags 		ServiceProviders
//
// @Security ApiKeyAuth
//
// @Accept  	json
// @Produce  	json
//
// @Param 		object	body 	object{service=models.Service} 	true	"Service Object"
//
// @Success 200 {object} models.BasicRes
//
// @Router /serviceproviders/add-service [post]
func AddServiceHandler(c *gin.Context, db *models.MongoDB) {
	var request struct {
		Service models.Service `json:"service"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	email := current_svcp.SVCPEmail
	err = svcp_utills.AddService(db, email, request.Service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Service added successfully"})
}

// DeleteBankAccountHandler godoc
//
// @Summary 	Delete bank account
// @Description Delete bank account (authentication required and only the service provider can update their own profile)
// @Tags 		ServiceProviders
//
// @Security ApiKeyAuth
//
// @Accept  	json
// @Produce  	json
//
// @Success 200 {object} models.BasicRes
//
// @Router /serviceproviders/delete-bank-account [delete]
func DeleteBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	err = svcp_utills.DeleteBankAccount(db, current_svcp.SVCPEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Bank account deleted successfully"})
}

// SetDefaultBankAccountHandler godoc
//
// @Summary 	Set default bank account
// @Description Set default bank account (authentication required and only the service provider can update their own profile)
// @Tags 		ServiceProviders
//
// @Security ApiKeyAuth
//
// @Accept  	json
// @Produce  	json
//
// @Param 		object	body	defaultBankAccountReq	true	"Default Bank Account Object"
//
// @Success 200 {object} models.BasicRes
//
// @Router /serviceproviders/set-default-bank-account [post]
func SetDefaultBankAccountHandler(c *gin.Context, db *models.MongoDB) {
	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	var request defaultBankAccountReq

	err = json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid request"})
		return
	}

	svcp_email := current_svcp.SVCPEmail
	default_bank_account := request.DefaultAccountNumber
	default_bank := request.DefaultBank

	_, err = svcp_utills.SetDefaultBankAccount(svcp_email, default_bank_account, default_bank, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Default bank account set successfully"})
}

type defaultBankAccountReq struct {
	DefaultAccountNumber string `json:"defaultAccountNumber"`
	DefaultBank          string `json:"defaultBank"`
}

// GetChatsHandler godoc
//
// @Summary 	Get the current service provider's chats
// @Description Get chat rooms of the current service provider. Each `Chat.messages` contains only *one* latest message.
// @Tags 		ServiceProviders
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
// @Router /serviceproviders/chats [get]
func GetChatsHandler(c *gin.Context, db *models.MongoDB) {
	// get current svcp
	current_svcp, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
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

	chats, err := chathistory.GetChatsById(db, current_svcp.SVCPID, page, per, "svcp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, chats)
}

func _authenticate(c *gin.Context, db *models.MongoDB) (*models.SVCP, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.User:
		err = errors.New("need token of type SVCP but recives token SVCP type")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, nil
		// Handle user
	case *models.SVCP:
		return entity, nil
		// Handle svcp
	}
	err = errors.New("need token of type SVCP but wrong type")
	c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
	return nil, err
}
