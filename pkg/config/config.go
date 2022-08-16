package config

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DB_USER = "postgres"
	DB_PASS = "postgres"
	DB_NAME = "postgres"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func GetTime() time.Time {
	return time.Now()
}

func GetCurrentTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(GetTime())
}

// Validate Data from request body
func ValidateData(data interface{}) bool {
	return data != nil
}

// Log Messages
func Logs(level string, message string) {
	switch level {
	case "info":
		log.Printf("INFO: %s --> %s", time.Now(), message)
	case "error":
		log.Printf("ERROR: %s --> %s", time.Now(), message)
	default:
		log.Printf("INFO: %s --> %s", time.Now(), message)
	}
}
