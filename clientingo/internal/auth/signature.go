package auth

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateSignature(url, token, timestamp string) string {
	data := url + token + timestamp
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
