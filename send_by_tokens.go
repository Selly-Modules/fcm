package fcm

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/thoas/go-funk"
)

// SendByTokens ...
func (s Service) SendByTokens(tokens []string, batchID string, payload messaging.Message) (result Result, err error) {
	ctx := context.Background()
	result.BatchID = batchID

	// Send
	for {
		listTokens, restTokens := separateTokens(tokens)

		// Return if list tokens length < 0
		if len(listTokens) <= 0 {
			break
		}

		// Prepare message
		message := &messaging.MulticastMessage{
			Tokens:       listTokens,
			Data:         payload.Data,
			Notification: payload.Notification,
			Android:      payload.Android,
		}

		// Send
		resp, e := s.Client.SendMulticast(ctx, message)
		if e != nil {
			err = e
			fmt.Printf("*** Error when push notification with batchID %s, error: %s \n", batchID, err.Error())
			return
		}

		result.Success += resp.SuccessCount
		result.Failure += resp.FailureCount

		// Assign token with rest tokens
		tokens = restTokens
	}

	// Save log
	go func() {
		if batchID != "" {
			log := Log{
				Action:       LogActionSendByTokens,
				BatchID:      batchID,
				Topics:       nil,
				TokenCount:   len(tokens),
				SuccessCount: result.Success,
				FailureCount: result.Failure,
			}
			s.saveLog(log)
		}
	}()

	return
}

// separate tokens for multiple times send, due to FCM limited
func separateTokens(tokens []string) (list, rest []string) {
	list = tokens
	if len(tokens) > maxTokensPerSend {
		list = removeEmptyToken(tokens[:maxTokensPerSend])
		rest = tokens[maxTokensPerSend:]
	}
	return
}

// remove empty token
func removeEmptyToken(tokens []string) []string {
	result := funk.FilterString(tokens, func(t string) bool {
		return t != ""
	})
	return result
}
