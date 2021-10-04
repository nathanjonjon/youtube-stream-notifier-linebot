package controller

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var Bot *linebot.Client

func InitLineBot() {
	// create line bot
	var err error
	Bot, err = linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Linebot created:", Bot)
}
