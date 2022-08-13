package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Rhaqim/thecommune-gobackend/config"
)

func SetupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DB_USER, config.DB_PASS, config.DB_NAME)

	db, err := sql.Open("postgres", dbinfo)

	log.Fatal(err)

	return db
}
