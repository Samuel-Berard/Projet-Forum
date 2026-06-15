package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func HashPasswordSHA512(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CheckPasswordHash(password, hash string) bool {
	return HashPasswordSHA512(password) == hash
}
