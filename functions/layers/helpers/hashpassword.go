package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

//GenerateHashPassword generate hashed password
func GenerateHashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

//CheckPasswordHash compare the none-hashed password with hashed string stored in database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
