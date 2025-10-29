package weather

import "fmt"

func (s *Service) FormatCurrentWeather(weather *WeatherResponse, city string) string {
	desc := s.getWeatherDescription(weather.Current.Weathercode)

	return fmt.Sprintf(
		"🌤️ *Погода в %s*\n\n"+
			"🌡️ Температура: *%.1f°C*\n"+
			"💨 Ветер: *%.1f км/ч*\n"+
			"💧 Влажность: *%d%%*\n"+
			"🌧️ Осадки: *%.1f мм*\n"+
			"📊 Состояние: *%s*\n\n"+
			"_Обновлено: %s_",
		city,
		weather.Current.Temperature,
		weather.Current.Windspeed,
		weather.Current.Humidity,
		weather.Current.Precipitation,
		desc,
		weather.Current.Time[:16],
	)
}

func (s *Service) getWeatherDescription(code int) string {
	weatherCodes := map[int]string{
		0: "☀️ Ясно", 1: "🌤️ Преимущественно ясно", 2: "⛅ Переменная облачность",
		3: "☁️ Пасмурно", 45: "🌫️ Туман", 48: "🌫️ Инейный туман",
		51: "🌧️ Легкая морось", 53: "🌧️ Умеренная морось", 55: "🌧️ Сильная морось",
		61: "🌧️ Небольшой дождь", 63: "🌧️ Умеренный дождь", 65: "🌧️ Сильный дождь",
		80: "🌦️ Небольшой ливень", 81: "🌦️ Умеренный ливень", 82: "🌦️ Сильный ливень",
		95: "⛈️ Гроза", 96: "⛈️ Гроза с градом",
	}

	if desc, exists := weatherCodes[code]; exists {
		return desc
	}
	return "❓ Неизвестно"
}
