package fcm

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
)

// SubscribeTokensToTopic ...
func (s Service) SubscribeTokensToTopic(batchID, topic string, tokens []string) (result Result, err error) {
	if topic == "" || len(tokens) == 0 {
		return
	}

	if !isTopicAllowed(topic) {
		return
	}

	ctx := context.Background()
	for {
		listTokens, restTokens := separateTokens(tokens)

		// Return if list tokens length < 0
		if len(listTokens) <= 0 {
			break
		}

		resp, e := s.Client.SubscribeToTopic(ctx, tokens, topic)
		if e != nil {
			err = e
			fmt.Printf("*** Subscribe tokens to topic %s error: %s \n", topic, err.Error())
			return
		}

		// Assign result
		result.Success += resp.SuccessCount
		result.Failure += resp.FailureCount

		// Get list error tokens
		if len(resp.Errors) > 0 {
			result.ErrorTokens = append(result.ErrorTokens, getFailureTokensFromSubscribe(resp.Errors, listTokens)...)
		}

		// Assign token with rest tokens
		tokens = restTokens
	}

	// Save log
	go func() {
		if batchID != "" {
			log := Log{
				Action:       LogActionSubscribeTokens,
				BatchID:      batchID,
				Topics:       []string{topic},
				TokenCount:   len(tokens),
				SuccessCount: result.Success,
				FailureCount: result.Failure,
			}
			s.saveLog(log)
		}
	}()

	return
}

func getFailureTokensFromSubscribe(r []*messaging.ErrorInfo, inputTokens []string) (tokens []string) {
	for _, info := range r {
		tokens = append(tokens, inputTokens[info.Index])
	}
	return
}
