package fcm

import (
	"context"
	"time"

	"github.com/Selly-Modules/mongodb"
	"github.com/kr/pretty"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log actions
const (
	LogActionSendByTokens      = "send_by_tokens"
	LogActionSendByTopic       = "send_by_topic"
	LogActionSubscribeTokens   = "subscribe_tokens"
	LogActionUnsubscribeTokens = "unsubscribe_tokens"
)

// Log ...
type Log struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Action       string             `bson:"action" json:"action"`
	BatchID      string             `bson:"batchId" json:"batchId"`
	Topics       []string           `bson:"topics" json:"topics"`
	TokenCount   int                `bson:"tokenCount" json:"tokenCount"`
	SuccessCount int                `bson:"successCount" json:"successCount"`
	FailureCount int                `bson:"failureCount" json:"failureCount"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

// Save log to db
func (s Service) saveLog(doc Log) {
	ctx := context.Background()

	// Assign data
	doc.ID = mongodb.NewObjectID()
	doc.CreatedAt = Now()

	if _, err := s.DB.Collection(logDBName).InsertOne(ctx, doc); err != nil {
		pretty.Println("*** FCM create log error", err.Error())
		pretty.Println("*** Payload", doc)
		pretty.Println("*****************")
	}
}
