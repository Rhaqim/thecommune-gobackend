package config

import "fmt"

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
