package main

import (
	"log"
	"os"

	"github.com/Tembo08/my-weather-bot/internal/bot"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	weatherBot := bot.NewBot(token)
	weatherBot.Run()
}
