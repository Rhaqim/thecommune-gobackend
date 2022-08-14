package main

/*
Handle password salting and hashing
*/

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TestPassword(password string) (string, error) {
	bytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(bytes[:]), nil
}

func main() {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hash)
	log.Println(CheckPasswordHash(password, hash))

	hash, err = TestPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hash)

}
