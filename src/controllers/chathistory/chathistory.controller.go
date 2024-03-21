package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/chat/chathistory"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetChatHistoryHandler godoc
//
// @Summary 	Get chat history
// @Description Get chat history of a room by roomId
// @Tags 		Chat
//
// @Accept  	json
// @Produce  	json
//
// @Param 		page	query	int 	false	"Page number(default 1)"
// @Param 		per 	query	int 	false 	"Number of items per page(default 10)"
//
// @Success 200 {object} models.Chat
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /chat/history/:roomId [get]
func GetChatHistoryHandler(c *gin.Context, db *models.MongoDB) {
	roomId := c.Param("roomId")
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

	chatHistory, err := chathistory.GetChatHistory(db, roomId, page, per)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chatHistory)
}

// CreateChatHistoryHandler godoc
//
// @Summary 	Create chat history
// @Description Create chat history of a room by roomId , user0Id , user1Id
// @Tags 		Chat
//
// @Accept  	json
// @Produce  	json
//
// @Success 200 {object} models.Chat
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /chat/history [post]
func CreateChatHistoryHandler(c *gin.Context, db *models.MongoDB) {
	type CreateChatHistory struct {
		RoomID    string `json:"roomId"`
		User0ID   string `json:"user0Id"`
		User1ID   string `json:"user1Id"`
		User0Type string `json:"user0Type"`
		User1Type string `json:"user1Type"`
	}
	var createChatHistoryReq CreateChatHistory
	err := c.BindJSON(&createChatHistoryReq)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	chatHistory, err := chathistory.CreateChatHistory(db, createChatHistoryReq.RoomID, createChatHistoryReq.User0ID, createChatHistoryReq.User1ID, createChatHistoryReq.User0Type, createChatHistoryReq.User1Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chatHistory)
}

// UpdateChatHistoryHandler godoc
//
// @Summary 	Update chat history
// @Description Update chat history of a room by roomId
// @Tags 		Chat
//
// @Accept  	json
// @Produce  	json
//
// @Success 200 {object} models.Chat
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /chat/history/:roomId [put]
func UpdateChatHistoryHandler(c *gin.Context, db *models.MongoDB) {
	roomId := c.Param("roomId")
	var updateChatHistoryReq bson.M
	err := c.BindJSON(&updateChatHistoryReq)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	updateChatHistory, err := chathistory.UpdateChatHistoryHandler(db, roomId, updateChatHistoryReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updateChatHistory)
}
