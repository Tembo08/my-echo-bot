package weather

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Time          string  `json:"time"`
		Temperature   float64 `json:"temperature_2m"`
		Windspeed     float64 `json:"windspeed_10m"`
		Weathercode   int     `json:"weathercode"`
		Precipitation float64 `json:"precipitation"`
		Humidity      int     `json:"relative_humidity_2m"`
	} `json:"current"`
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		Weathercode    []int     `json:"weathercode"`
	} `json:"daily"`
}
