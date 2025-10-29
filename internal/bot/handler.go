package bot

import (
	"my-weather-bot/internal/location"
	"my-weather-bot/internal/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	weatherService  *weather.Service
	locationService *location.Service
}

func NewHandler() *Handler {
	return &Handler{
		weatherService:  weather.NewService(),
		locationService: location.NewService(),
	}
}

func (h *Handler) HandleStart(bot *tgbotapi.BotAPI, chatID int64) {
	message := `🌤️ *Бот-метеоролог*

Я могу показать актуальную погоду!

*Команды:*
/weather - погода в Москве
/city <город> - погода для конкретного города
/help - все команды

Пример: /city Санкт-Петербург`

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *Handler) HandleHelp(bot *tgbotapi.BotAPI, chatID int64) {
	message := `*Доступные команды:*

/weather - текущая погода (Москва)
/city <город> - погода в указанном городе
/help - эта справка

*Примеры:*
/city Лондон
/city Париж
/city Новосибирск`

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *Handler) HandleWeather(bot *tgbotapi.BotAPI, chatID int64) {
	// По умолчанию Москва
	lat, lon := 55.7558, 37.6173
	city := "Москве"

	weatherData, err := h.weatherService.GetCurrentWeather(lat, lon)
	if err != nil {
		h.sendError(bot, chatID, "Ошибка получения погоды")
		return
	}

	message := h.weatherService.FormatCurrentWeather(weatherData, city)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *Handler) HandleCityWeather(bot *tgbotapi.BotAPI, chatID int64, cityName string) {
	if cityName == "" {
		msg := tgbotapi.NewMessage(chatID, "Укажите город: `/city Москва`")
		msg.ParseMode = "Markdown"
		bot.Send(msg)
		return
	}

	coords, err := h.locationService.GetCoordinates(cityName)
	if err != nil {
		h.sendError(bot, chatID, "Город не найден")
		return
	}

	weatherData, err := h.weatherService.GetCurrentWeather(coords.Latitude, coords.Longitude)
	if err != nil {
		h.sendError(bot, chatID, "Ошибка получения погоды")
		return
	}

	message := h.weatherService.FormatCurrentWeather(weatherData, coords.Name)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *Handler) sendError(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	bot.Send(msg)
}
