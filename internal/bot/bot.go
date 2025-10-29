package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WeatherBot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

func NewBot(token string) *WeatherBot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &WeatherBot{
		api:     api,
		handler: NewHandler(),
	}
}

func (b *WeatherBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update)
		}
	}
}

func (b *WeatherBot) handleCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		b.handler.HandleStart(b.api, chatID)
	case "weather":
		b.handler.HandleWeather(b.api, chatID)
	case "city":
		city := update.Message.CommandArguments()
		b.handler.HandleCityWeather(b.api, chatID, city)
	case "help":
		b.handler.HandleHelp(b.api, chatID)
	}
}
