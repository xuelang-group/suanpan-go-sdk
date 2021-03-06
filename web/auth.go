package web

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"

	"github.com/xuelang-group/suanpan-go-sdk/config"
)

func signatureV1(secret, data string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(data))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func GetHeaders() map[string][]string {
	env := config.GetEnv()
	headers := make(map[string][]string)
	headers[env.SpUserIdHeaderField] = []string{env.SpUserId}
	headers[env.SpUserSignatureHeaderField] = []string{signatureV1(env.SpAccessSecret, env.SpUserId)}
	headers[env.SpUserSignVersionHeaderField] = []string{"v1"}

	return headers
}