package token

import (
	"crypto/rand"
	"encoding/base64"
)

const size = 32

func GenerateToken() (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
