package controllers

import (
	"errors"
	"io"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	issue_utills "petpal-backend/src/utills/issue"

	"github.com/gin-gonic/gin"
)

// CreateIssue godoc
//
// @Summary     Create issue
// @Description Create an issue for the current user
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
// @Success     202      {object} models.BasicRes    "Accepted"
// @Failure     400      {object} models.BasicErrorRes      "Bad request"
// @Failure     401      {object} models.BasicErrorRes      "Unauthorized"
// @Failure     500      {object} models.BasicErrorRes      "Internal server error"
//
// @Router      /user/uploadProfileImage [post]
func CreateIssue(c *gin.Context, db *models.MongoDB) {
	user, svcp, err := _authenticate(c, db)
	if err != nil {
		return
	}

	issue := models.CreateIssue{}

	file, _, err := c.Request.FormFile("file")
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

	if user != nil {
		issue.ReporterID = user.ID
		issue.ReporterType = "user"
	} else if svcp != nil {
		issue.ReporterID = svcp.SVCPID
		issue.ReporterType = "svcp"
	} else {
		err = errors.New("error getting user or svcp object from token")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
	}

	if err := issue_utills.CreateIssue(db, &issue); err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BasicRes{Message: "Issue created successfully"})
}

func _authenticate(c *gin.Context, db *models.MongoDB) (*models.User, *models.SVCP, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, nil, err
	}
	switch entity := entity.(type) {
	case *models.User:
		return entity, nil, nil
		// Handle user
	case *models.SVCP:
		return nil, entity, nil
		// Handle svcp
	default:
		err = errors.New("need token of type User or SVCP but received Admin token")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, nil, err
	}
}
