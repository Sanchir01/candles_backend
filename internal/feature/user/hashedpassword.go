package user

import (
	"crypto/subtle"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return passHash, nil
}

func VerifyPassword(password, hash string) bool {
	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(hash)) == 1
}
