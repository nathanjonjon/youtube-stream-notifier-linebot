package controller

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var DB *sql.DB

func InitDatabase() {
	// connect to the db
	DB, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	log.Println("DB conneted:", DB)
}
