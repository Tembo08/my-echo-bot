package main

import (
	"log"
	"os"

	"my-weather-bot/internal/bot"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env из корня проекта
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	weatherBot := bot.NewBot(token)
	weatherBot.Run()
}
