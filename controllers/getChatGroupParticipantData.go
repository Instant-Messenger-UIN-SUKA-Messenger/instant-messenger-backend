package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetChatGroupParticipantData(c *gin.Context) {
	participantsId := c.QueryArray("participantsId")

	if len(participantsId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Missing required parameter: participantsId",
		})
		return
	}

	log.Printf("Extracted participantsId: %s\n", participantsId)

	// // Get Participant Data based on participantsId
	// for _, participantId := range participantsId {
	// 	// Find chat information based on chatID
	// 	curChat, err := chatCollections.Find(context.TODO(), bson.M{"chatId": chatId})
	// 	if err != nil {
	// 		log.Printf("Error finding chat with ID %s: %v\n", chatId, err)
	// 	}

	// 	defer curChat.Close(context.TODO())
	// 	for curChat.Next(context.TODO()) {
	// 		var chat models.Chat
	// 		err := curChat.Decode(&chat)
	// 		log.Printf("%s", chat)
	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, gin.H{
	// 				"error":   true,
	// 				"message": err.Error(),
	// 			})
	// 			return
	// 		}

	// 		// Check if chat has a name (because if private chat, chat name is empty)
	// 		if chat.Name == "" {
	// 			// Get Other Chat Participant userId
	// 			var otherParticipantID string
	// 			if len(chat.Participants) == 2 { // Assuming private chat has exactly 2 participants
	// 				// Extract the other participant's userId based on whether it matches the current user's nim
	// 				if chat.Participants[0] == nim {
	// 					otherParticipantID = chat.Participants[1]
	// 				} else {
	// 					otherParticipantID = chat.Participants[0]
	// 				}
	// 			}

	// 			// Fetch user details of the other participant based on their userId/nim
	// 			var otherParticipantUser models.User
	// 			err := userCollections.FindOne(context.TODO(), bson.M{"nim": otherParticipantID}).Decode(&otherParticipantUser)
	// 			if err != nil {
	// 				// Handle error fetching participant details
	// 				log.Printf("Error finding participant for chat %s: %v\n", chat.ChatID, err)
	// 				// You can decide how to handle this error (e.g., continue or return error)
	// 				continue
	// 			}

	// 			// Set chat name to other participant's name
	// 			chat.Name = otherParticipantUser.Name
	// 		}

	// 		chats = append(chats, chat)
	// 	}
	// }
}
