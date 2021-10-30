package authentification

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRememberToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
