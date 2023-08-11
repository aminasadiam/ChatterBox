package generator

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateActivationCode() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}
