package fcm

import (
	"context"
	"fmt"
	"strings"

	"firebase.google.com/go/messaging"
)

// SendByTopics ...
func (s Service) SendByTopics(topics []string, batchID string, payload messaging.Message) {
	ctx := context.Background()

	// Return if have no topics
	if len(topics) == 0 {
		fmt.Sprintf("*** Empty topics array with batch id %s \n", batchID)
		return
	}

	// Get topic condition
	payload.Condition = getTopicCondition(topics)

	// Return if there is no condition
	if payload.Condition == "" {
		fmt.Sprintf("*** No valid topics array with batch id %s: %v \n", batchID, topics)
		return
	}

	_, err := s.Client.Send(ctx, &payload)
	if err != nil {
		fmt.Sprintf("*** Send topic error with batch id %s: %s \n", batchID, err.Error())
		fmt.Sprintf("*** Topics: %v \n", topics)
	}

	// Save log
	go func() {
		if batchID != "" {
			log := Log{
				Action:       LogActionSendByTopic,
				BatchID:      batchID,
				Topics:       topics,
				TokenCount:   0,
				SuccessCount: 0,
				FailureCount: 0,
			}
			s.saveLog(log)
		}
	}()
}

// getTopicCondition ...
func getTopicCondition(topics []string) string {
	var conditions []string

	for _, topic := range topics {
		if !isTopicAllowed(topic) {
			continue
		}
		cond := fmt.Sprintf("'%s' in topics", topic)
		conditions = append(conditions, cond)
	}

	if len(conditions) == 0 {
		return ""
	}
	return strings.Join(conditions, " && ")
}
