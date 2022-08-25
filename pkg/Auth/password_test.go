package auth

/*
Handle password salting and hashing
*/

import (
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func _HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func _CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func encoding(password string) (string, error) {
// 	bytes := sha256.Sum256([]byte(password))
// 	return hex.EncodeToString(bytes[:]), nil
// }

func TestPassword(t *testing.T) {
	password := "password"
	hash, err := _HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hash)
	log.Println(_CheckPasswordHash(password, hash))
}
