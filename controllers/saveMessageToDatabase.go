package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"instant-messenger-backend/database"
	"instant-messenger-backend/models"
)

func SaveMessageToDatabase(messageJSON []byte) error {
    // Decode JSON message
    var message models.Message
    err := json.Unmarshal(messageJSON, &message)
    if err != nil {
            return fmt.Errorf("error unmarshalling message: %w", err)
    }

    // Generate a new ObjectID
    message.ID = primitive.NewObjectID()

    // Insert the message into the database
    collection := database.GetCollection(database.DB, "messages")
    _, err = collection.InsertOne(context.TODO(), message)
    if err != nil {
            return fmt.Errorf("error inserting message into database: %w", err)
    }

    return nil // No need to return the message here
}