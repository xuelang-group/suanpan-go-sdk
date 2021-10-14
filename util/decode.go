package util

import (
	"encoding/base64"
)

func DecodeBase64(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(b), nil
}