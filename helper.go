package fcm

import (
	b64 "encoding/base64"

	"github.com/thoas/go-funk"
)

// base64Decode ...
func base64Decode(text string) []byte {
	sDec, _ := b64.StdEncoding.DecodeString(text)
	return sDec
}

// isTopicAllowed ...
func isTopicAllowed(topic string) bool {
	return funk.ContainsString(allowedTopics, topic)
}
