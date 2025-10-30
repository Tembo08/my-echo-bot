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
	log.Println("🔧 Initializing new WeatherBot...")

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("❌ Failed to create BotAPI:", err)
	}

	api.Debug = true
	log.Printf("✅Authorized on account %s", api.Self.UserName)

	return &WeatherBot{
		api:     api,
		handler: NewHandler(),
	}
}

func (b *WeatherBot) Run() {
	log.Println("🚀 Starting bot main loop...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	log.Println("📡 Bot is listening for updates...")

	for update := range updates {
		log.Printf("📨 Received update ID: %d", update.UpdateID)

		if update.Message == nil {
			log.Println("⏭️  Skip: message is nil")
			continue
		}

		log.Printf("👤 Message from: %s (ID: %d)", update.Message.From.UserName, update.Message.From.ID)
		log.Printf("💬 Text: %s", update.Message.Text)

		if update.Message.IsCommand() {
			log.Println("🔍 Detected command")
			b.handleCommand(update)
		} else {
			log.Println("💭 Regular message (not a command)")
		}

	}
}

func (b *WeatherBot) handleCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	command := update.Message.Command()

	log.Printf("🎯 Handling command: /%s", command) // ← Добавил логирование команды

	switch command {

	case "start":
		b.handler.HandleStart(b.api, chatID)
		log.Println("✅ You pressed Start!")

	case "weather":
		b.handler.HandleWeather(b.api, chatID)
		log.Println("✅ You pressed Weather!")

	case "help":
		b.handler.HandleHelp(b.api, chatID)
		log.Println("✅ You pressed Help!")

	default:
		log.Printf("⚠️  Unknown command: /%s", command) // ← Уточнил какая команда
	}

}
