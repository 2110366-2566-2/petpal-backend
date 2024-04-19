package controllers

import (
	"errors"
	"io"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	issue_utills "petpal-backend/src/utills/issue"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateIssue godoc
//
// @Summary     Create issue
// @Description Create an issue for the current user. If attached image is not provided, it will be set to null.
// @Tags        Issue
//
// @Accept      multipart/form-data
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       attachedImage      	formData    file      	false		"Attached image (optional)"
// @Param       details            	formData    string      true		"Details of issue"
// @Param       issueType      		formData    string      true        "Type of issue (refund, system, service)"
// @Param       associatedBookingID	formData    string      false 		"ID of associated booking if type is service (optional)"
//
// @Success     200      {object} models.BasicRes    "Accepted"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /issue [post]
func CreateIssue(c *gin.Context, db *models.MongoDB) {
	e, e_type, err := _authenticate(c, db)
	if err != nil {
		return
	}

	issue := models.CreateIssue{}

	file, _, err := c.Request.FormFile("attachedImage")
	if err != nil && err.Error() != "http: no such file" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if err == nil {
		attachedImg, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Error reading file: " + err.Error()})
			return
		}
		issue.AttachedImg = attachedImg
		defer file.Close()
	}

	issue.Details = c.Request.FormValue("details")
	issue.IssueType = c.Request.FormValue("issueType")
	issue.AssociatedBookingID = c.Request.FormValue("associatedBookingID")

	if e_type == "user" {
		issue.ReporterID = e.(*models.User).ID
		issue.ReporterType = "user"
	} else if e_type == "svcp" {
		issue.ReporterID = e.(*models.SVCP).SVCPID
		issue.ReporterType = "svcp"
	} else {
		err = errors.New("error getting user or svcp object from token")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if err := issue_utills.CreateIssue(db, &issue); err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Issue created successfully"})
}

// GetIssue godoc
//
// @Summary     Get issues
// @Description Get issues associated with the current user. If user is admin, all issues are returned. Otherwise only issues associated with the user are returned. Note that in an issue, if attached image is not provided, it will be set to null.
// @Tags        Issue
//
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       page    query    int     false        "Page number"
// @Param       per     query    int     false        "Number of issues per page"
//
// @Success     200      {object} []models.Issue    			"Success"
// @Failure     400      {object} models.BasicErrorRes  		"Bad request"
// @Failure     401      {object} models.BasicErrorRes  		"Unauthorized"
// @Failure     500      {object} models.BasicErrorRes  		"Internal server error"
//
// @Router      /issue [get]
func GetIssues(c *gin.Context, db *models.MongoDB) {
	e, e_type, err := _authenticate(c, db)
	if err != nil {
		return
	}

	// set default page and per
	params := c.Request.URL.Query()
	if !params.Has("page") {
		params.Set("page", "1")
	}
	if !params.Has("per") {
		params.Set("per", "10")
	}

	// parse page and per
	page, err_page := strconv.ParseInt(params.Get("page"), 10, 64)
	per, err_per := strconv.ParseInt(params.Get("per"), 10, 64)

	if err_page != nil || err_per != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid page or per"})
		return
	}

	filter := bson.M{}
	if e_type == "user" {
		filter["reporterID"] = e.(*models.User).ID
		filter["reporterType"] = "user"
	} else if e_type == "svcp" {
		filter["reporterID"] = e.(*models.SVCP).SVCPID
		filter["reporterType"] = "svcp"
	} else if e_type == "admin" {
		// do nothing
	} else {
		err = errors.New("error getting user or svcp object from token")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	// get all issues
	issues, err := issue_utills.GetIssues(db, filter, page-1, per)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, issues)
}

// GetIssueById godoc
//
// @Summary     Get issue by ID
// @Description Get issue by ID
// @Tags        Issue
//
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       id      path    string     true        "ID of issue to get"
//
// @Success     200      {object} models.Issue    		"Success"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /issue/{id} [get]
func GetIssueById(c *gin.Context, db *models.MongoDB) {
	issueID := c.Param("id")

	issue, err := issue_utills.GetIssueByID(db, issueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// AdminAcceptIssue godoc
//
// @Summary     Admin accept issue
// @Description Admin accepts an issue. This will set the workingAdminID field of the issue to the admin's ID. If the issue is already accepted by another admin, this will return an error.
// @Tags        Issue
//
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       id      path    string     true        "ID of issue to accept"
//
// @Success     200      {object} models.BasicRes    		"Accepted"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /issue/accept/{id} [post]
func AdminAcceptIssue(c *gin.Context, db *models.MongoDB) {
	e, e_type, err := _authenticate(c, db)
	if err != nil {
		return
	}

	if e_type != "admin" {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: "Only admin can accept issues"})
		return
	}

	admin := e.(*models.Admin)

	issueID := c.Param("id")

	if err := issue_utills.AdminAcceptIssue(db, issueID, admin.AdminID); err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Issue accepted by " + admin.AdminID})

}

// AdminResolveIssue godoc
//
// @Summary     Admin resolve issue
// @Description Admin resolves an issue. This will set the isResolved field of the issue to true and the resolveDate field to the current date. The issue can only be resolved by the admin who accepted the issue. If the current admin is not the working admin of the issue, this will return an error saying that no issue is found.
// @Tags        Issue
//
// @Produce     json
//
// @Security    ApiKeyAuth
//
// @Param       id      path    string     true        "ID of issue to resolve"
//
// @Success     200      {object} models.BasicRes    		"Accepted"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /issue/resolve/{id} [post]
func AdminResolveIssue(c *gin.Context, db *models.MongoDB) {
	e, e_type, err := _authenticate(c, db)
	if err != nil {
		return
	}

	if e_type != "admin" {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: "Only admin can resolve issues"})
		return
	}

	admin := e.(*models.Admin)

	issueID := c.Param("id")

	if err := issue_utills.AdminResolveIssue(db, issueID, admin.AdminID); err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Issue resolved by " + admin.AdminID})
}

func _authenticate(c *gin.Context, db *models.MongoDB) (interface{}, string, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, "", err
	}
	switch entity := entity.(type) {
	case *models.User:
		return entity, "user", nil
		// Handle user
	case *models.SVCP:
		return entity, "svcp", nil
		// Handle svcp
	case *models.Admin:
		return entity, "admin", nil
	default:
		err = errors.New("unexpected login type")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, "", err
	}
}
