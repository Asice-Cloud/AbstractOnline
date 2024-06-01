package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// lowercase
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tmpStr := h.Sum(nil)

	return hex.EncodeToString(tmpStr)
}

// uppercase
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// crypt
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// decrypt
func ValidPassword(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
