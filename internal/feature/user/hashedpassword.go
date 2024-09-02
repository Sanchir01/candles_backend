package userFeature

import (
	"crypto/subtle"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

func GeneratePasswordHash(password string) string {
	salt := make([]byte, 16)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func VerifyPassword(password, hash string) bool {
	hashedPassword := GeneratePasswordHash(password)
	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(hash)) == 1
}
