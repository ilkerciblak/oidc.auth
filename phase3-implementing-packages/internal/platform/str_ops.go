package platform

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)


func GenerateBase64EncodedString() (string, error) {
	randomBytes := make([]byte, 32)

	// generated random bytes using crypto/rand
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("[rand.Read]:  %v", err)
	}

	state := base64.URLEncoding.EncodeToString(randomBytes)
	return state, nil
}
