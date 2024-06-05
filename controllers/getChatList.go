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
)

var userCollections *mongo.Collection = database.GetCollection(database.DB, "users")
var chatCollections *mongo.Collection = database.GetCollection(database.DB, "chats")

func GetChatList(c *gin.Context) {
	// Get data chat list based on nim from user collection
	var user models.User
	var nim string
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}
	nim = user.NIM
	log.Printf("Extracted nim: %s\n", nim)

	var chatList []string

	curUser, err := userCollections.Find(context.TODO(), bson.M{"nim": nim})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	defer curUser.Close(context.TODO())
	for curUser.Next(context.TODO()) {
		var user models.User
		err := curUser.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		chatList = user.ChatList
	}

	// Get data chat information based on chatList from chat collection
	var chats []models.Chat
	for _, chatId := range chatList {
		// Find chat information based on chatID
		curChat, err := chatCollections.Find(context.TODO(), bson.M{"chatId": chatId})
		if err != nil {
			log.Printf("Error finding chat with ID %s: %v\n", chatId, err)
		}

		defer curChat.Close(context.TODO())
		for curChat.Next(context.TODO()) {
			var chat models.Chat
			err := curChat.Decode(&chat)
			log.Printf("%s", chat)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": err.Error(),
				})
				return
			}
			chats = append(chats, chat)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"chats": chats,
	})
}
