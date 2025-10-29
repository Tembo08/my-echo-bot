package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Lat      float64         `json:"lat"`
	Lon      float64         `json:"lon"`
	Timezone string          `json:"timezone"`
	Daily    []DailyForecast `json:"daily"`
}

type DailyForecast struct {
	Dt   int64 `json:"dt"`
	Temp struct {
		Day   float64 `json:"day"`
		Night float64 `json:"night"`
	} `json:"temp"`
	Weather   []WeatherInfo `json:"weather"`
	Humidity  int           `json:"humidity"`
	WindSpeed float64       `json:"wind_speed"`
	Rain      float64       `json:"rain,omitempty"`
	Snow      float64       `json:"snow,omitempty"`
}

type WeatherInfo struct {
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Структура для IP-геолокации
type IPGeoResponse struct {
	IP       string  `json:"ip"`
	City     string  `json:"city"`
	Region   string  `json:"region"`
	Country  string  `json:"country"`
	Loc      string  `json:"loc"`
	Lat      float64 `json:"lat,omitempty"`
	Lon      float64 `json:"lon,omitempty"`
	Timezone string  `json:"timezone"`
}

type GeoCodingResponse struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

func main() {

	// Загружаем env.
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Warning: env.file not found")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("TG_BOT_TOKEN not set")
	}

	// создаем экземпляр бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Bot is starting...")

	// Настраиваем канал обновлений

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем обновления от Телеги
	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(bot, update.Message)
			continue
		}

		//Если просто текст - предлагаем помощь
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"🌤️ Используй /weather для погоды в твоём городе\n/help - все команды")
		bot.Send(msg)

	}
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := message.Command()
	args := strings.TrimSpace(message.CommandArguments())

	switch command {
	case "start":
		sendWelcome(bot, message.Chat.ID)
	case "help":
		sendHelp(bot, message.Chat.ID)
	case "weather", "погода":
		getWeatherByLocation(bot, message)
	case "город", "city":
		if args == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID,
				"Укажи город: /город Москва")
			bot.Send(msg)
			return
		}

		getWeatherByCity(bot, message.Chat.ID, args)

	default:
		msg := tgbotapi.NewMessage(message.Chat.ID,
			"Неизвестная команда Исползуй /help")
		bot.Send(msg)
	}

}

func sendWelcome(bot *tgbotapi.BotAPI, chatID int64) {

	welcomeText := `🌤️ *Бот-метеоролог*

	Я могу показать прогноз погоды на неделю!

	*Команды:*
/погода - погода в твоём городе (определяется автоматически)
/город Москва - погода для конкретного города
/help - все команды

Просто нажми /погода и я покажу погоду там, где ты находишься!`

	msg := tgbotapi.NewMessage(chatID, welcomeText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

}

func sendHelp(bot *tgbotapi.BotAPI, chatID int64) {
	helpText := `🌤️ *Доступные команды:*

/weather - погода в твоём городе (автоопределение)
/погода - то же самое
/city Москва - погода для конкретного города
/help - эта справка

*Примеры:*
/weather - погода там где ты есть
/city London - погода в Лондоне
/city Берлин - погода в Берлине`

	msg := tgbotapi.NewMessage(chatID, helpText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func getWeatherByLocation(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Отправляем сообщение "Определяем местоположение..."
	msg := tgbotapi.NewMessage(chatID, "📍 Определяю твоё местоположение...")
	bot.Send(msg)

	// Получаем приблизительное местоположение по IP
	location, err := getLocationByIP()
	if err != nil {
		log.Printf("Location error: %v", err)
		msg := tgbotapi.NewMessage(chatID,
			"❌ Не удалось определить местоположение. Используй /city Город")
		bot.Send(msg)
		return
	}

	// Получаем прогноз погоды
	weather, err := getWeatherForecast(location.Lat, location.Lon)
	if err != nil {
		log.Printf("Weather API error: %v", err)
		msg := tgbotapi.NewMessage(chatID,
			"❌ Ошибка получения данных о погоде. Попробуй позже.")
		bot.Send(msg)
		return
	}

	// Формируем красивый ответ
	cityName := fmt.Sprintf("%s, %s", location.City, location.Country)
	response := formatWeatherResponse(cityName, weather, true)
	msg = tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func getWeatherByCity(bot *tgbotapi.BotAPI, chatID int64, city string) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🔍 Ищу погоду для: %s", city))
	bot.Send(msg)

	// Получаем координаты города
	coords, err := getCityCoordinates(city)
	if err != nil {
		log.Printf("Geocoding error: %v", err)
		msg := tgbotapi.NewMessage(chatID,
			"❌ Не удалось найти город. Проверь название.")
		bot.Send(msg)
		return
	}

	// Получаем прогноз погоды
	weather, err := getWeatherForecast(coords.Lat, coords.Lon)
	if err != nil {
		log.Printf("Weather API error: %v", err)
		msg := tgbotapi.NewMessage(chatID,
			"❌ Ошибка получения данных о погоде. Попробуй позже.")
		bot.Send(msg)
		return
	}

	// Формируем ответ
	cityName := fmt.Sprintf("%s, %s", coords.Name, coords.Country)
	response := formatWeatherResponse(cityName, weather, false)
	msg = tgbotapi.NewMessage(chatID, response)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// Получение местоположения по IP
func getLocationByIP() (*IPGeoResponse, error) {
	resp, err := http.Get("http://ipinfo.io/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location IPGeoResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	// Парсим координаты из формата "lat,lon"
	if location.Loc != "" {
		parts := strings.Split(location.Loc, ",")
		if len(parts) == 2 {
			location.Lat, _ = strconv.ParseFloat(parts[0], 64)
			location.Lon, _ = strconv.ParseFloat(parts[1], 64)
		}
	}

	return &location, nil
}

// Получение координат города
func getCityCoordinates(city string) (*GeoCodingResponse, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY not set")
	}

	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geocode GeoCodingResponse
	err = json.Unmarshal(body, &geocode)
	if err != nil {
		return nil, err
	}

	if len(geocode) == 0 {
		return nil, fmt.Errorf("город не найден")
	}

	return &geocode, nil
}

// Получение прогноза погоды
func getWeatherForecast(lat, lon float64) (*WeatherResponse, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY not set")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%f&lon=%f&exclude=current,minutely,hourly,alerts&appid=%s&units=metric&lang=ru", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

// Форматирование ответа с погодой
func formatWeatherResponse(cityName string, weather *WeatherResponse, isCurrentLocation bool) string {
	var builder strings.Builder

	locationType := "📍"
	if !isCurrentLocation {
		locationType = "🏙️"
	}

	builder.WriteString(fmt.Sprintf("%s *%s*\n\n", locationType, cityName))
	builder.WriteString("*Прогноз на неделю:*\n\n")

	// Показываем прогноз на 5 дней
	for i := 0; i < 5 && i < len(weather.Daily); i++ {
		day := weather.Daily[i]
		date := time.Unix(day.Dt, 0)
		weekday := getRussianWeekday(date.Weekday())

		// Эмодзи для погоды
		weatherEmoji := getWeatherEmoji(day.Weather[0].Main, day.Weather[0].Description)

		builder.WriteString(fmt.Sprintf("*%s, %s* %s\n", weekday, date.Format("02.01"), weatherEmoji))
		builder.WriteString(fmt.Sprintf("   🌡️ %.1f°C (ночью %.1f°C)\n", day.Temp.Day, day.Temp.Night))
		builder.WriteString(fmt.Sprintf("   💨 %.1f м/с\n", day.WindSpeed))
		builder.WriteString(fmt.Sprintf("   💧 %d%%\n", day.Humidity))

		if day.Rain > 0 {
			builder.WriteString(fmt.Sprintf("   🌧️ %.1f мм\n", day.Rain))
		}
		if day.Snow > 0 {
			builder.WriteString(fmt.Sprintf("   ❄️ %.1f мм\n", day.Snow))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// Получение русского названия дня недели
func getRussianWeekday(weekday time.Weekday) string {
	days := map[time.Weekday]string{
		time.Monday:    "Пн",
		time.Tuesday:   "Вт",
		time.Wednesday: "Ср",
		time.Thursday:  "Чт",
		time.Friday:    "Пт",
		time.Saturday:  "Сб",
		time.Sunday:    "Вс",
	}
	return days[weekday]
}

// Получение эмодзи для типа погоды
func getWeatherEmoji(main, description string) string {
	switch main {
	case "Clear":
		return "☀️"
	case "Clouds":
		if strings.Contains(description, "few") || strings.Contains(description, "scattered") {
			return "⛅"
		}
		return "☁️"
	case "Rain":
		if strings.Contains(description, "light") {
			return "🌦️"
		}
		return "🌧️"
	case "Drizzle":
		return "🌦️"
	case "Thunderstorm":
		return "⛈️"
	case "Snow":
		return "❄️"
	case "Mist", "Fog", "Haze":
		return "🌫️"
	default:
		return "🌤️"
	}
}
