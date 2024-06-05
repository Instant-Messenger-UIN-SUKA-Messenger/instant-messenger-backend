package controllers

import (
	"context"
	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var chatCollections *mongo.Collection = database.GetCollection(database.DB, "chats")
var messageCollections *mongo.Collection = database.GetCollection(database.DB, "messages")

func GetListChat(c *gin.Context) {
	var chats []models.Chat

	cur, err := chatCollections.Find(context.TODO(), bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var chat models.Chat
		err := cur.Decode(&chat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		chats = append(chats, chat)
	}

	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
	})
}

func GetChatsDetail(c *gin.Context) {
	chatID := c.Param("id")

	chatObjectID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid chat ID",
		})
		return
	}

	var chat models.Chat
	err = chatCollections.FindOne(context.TODO(), bson.M{"_id": chatObjectID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Chat not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	var messages []models.Message
	cur, err := messageCollections.Find(context.TODO(), bson.M{"chatId": chatID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var message models.Message
		err := cur.Decode(&message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		messages = append(messages, message)
	}

	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat":     chat,
		"messages": messages,
	})
}
