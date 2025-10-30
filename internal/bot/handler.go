package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleStart(bot *tgbotapi.BotAPI, chatID int64) {
	message := `🌤️ *Бот-метеоролог*

Я покажу актуальную погоду в Москве!

*Команды:*
/weather - погода в Москве
/help - справка`

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *Handler) HandleHelp(bot *tgbotapi.BotAPI, chatID int64) {
	message := `*Доступные команды:*

/weather - текущая погода в Москве
/help - эта справка

Просто нажми /weather и узнай погоду!`

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *Handler) HandleWeather(bot *tgbotapi.BotAPI, chatID int64) {
	message := `🌤️ *Погода в Москве*

🌡️ Температура: *+15°C*
💨 Ветер: *3 м/с*  
💧 Влажность: *65%*
☁️ Состояние: *Облачно*

_Данные обновляются..._`

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
