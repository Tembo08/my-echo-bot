package location

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type GeoResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name"`
		Country   string  `json:"country"`
	} `json:"results"`
}

func (s *Service) GetCoordinates(city string) (*Coordinates, error) {
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", city)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geo GeoResponse
	err = json.NewDecoder(resp.Body).Decode(&geo)
	if err != nil || len(geo.Results) == 0 {
		return nil, fmt.Errorf("город не найден")
	}

	result := geo.Results[0]
	return &Coordinates{
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Name:      result.Name,
		Country:   result.Country,
	}, nil
}
