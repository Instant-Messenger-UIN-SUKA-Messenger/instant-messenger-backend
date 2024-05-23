package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    ChatID    string             `bson:"chatID,omitempty" json:"chatID,omitempty"`
    SenderID  string             `bson:"senderID,omitempty" json:"senderID,omitempty"`
    Content   string             `bson:"content" json:"content"`
    SentAt    time.Time          `bson:"sentAt" json:"sentAt"`
    Status    string             `bson:"status" json:"status"`
    Attachments []string         `bson:"attachments,omitempty" json:"attachments,omitempty"`
}
