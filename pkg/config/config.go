package config

import (
	"fmt"
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
