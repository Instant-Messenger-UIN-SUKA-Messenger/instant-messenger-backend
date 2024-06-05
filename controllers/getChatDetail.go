package controllers

import (
	"context"
	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var messageCollections *mongo.Collection = database.GetCollection(database.DB, "messages")

func GetChatDetail(c *gin.Context) {
	// Get chatId from body request
	var chat models.Chat
	var chatId string
	if err := c.BindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true, 
			"message": "Invalid request body",
		})
		return
	}
	chatId = chat.ChatID
	log.Printf("Extracted chatId: %s\n", chatId)

	// Get All message with chatId matches And Sorting Asc Based SentAt field
	var messageList []models.Message
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"SentAt": 1})
	curMessage, err := messageCollections.Find(context.TODO(), bson.M{"chatId": chatId}, findOptions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer curMessage.Close(context.TODO())
	for curMessage.Next(context.TODO()) {
		var message models.Message
		err := curMessage.Decode(&message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		messageList = append(messageList, message)
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"messages": messageList,
	})
}
