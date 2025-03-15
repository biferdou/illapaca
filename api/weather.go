package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/biferdou/illapaca/config"
	"github.com/biferdou/illapaca/model"
	"github.com/briandowns/spinner"
)

const (
	// WeatherAPI.com base URL
	baseURL = "https://api.weatherapi.com/v1"
)

// FetchWeather retrieves weather data from the API
func FetchWeather(location string, days int) (*model.WeatherData, error) {
	if config.AppConfig.APIKey == "" {
		return nil, fmt.Errorf("API key not set. Use --api-key flag or set ILLAPA_API_KEY environment variable")
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Fetching weather data "
	s.Start()
	defer s.Stop()

	// Construct forecast URL
	forecastURL := fmt.Sprintf("%s/forecast.json?key=%s&q=%s&days=%d&aqi=no&alerts=no",
		baseURL, config.AppConfig.APIKey, location, days)

	// Make request
	resp, err := http.Get(forecastURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Convert to our unified model
	weatherData := &model.WeatherData{
		Current: model.CurrentWeather{
			TempC:      response.Current.TempC,
			TempF:      response.Current.TempF,
			IsDay:      response.Current.IsDay,
			Condition:  response.Current.Condition,
			WindMph:    response.Current.WindMph,
			WindKph:    response.Current.WindKph,
			WindDir:    response.Current.WindDir,
			PressureMb: response.Current.PressureMb,
			PrecipMm:   response.Current.PrecipMm,
			Humidity:   response.Current.Humidity,
			FeelsLikeC: response.Current.FeelslikeC,
			FeelsLikeF: response.Current.FeelslikeF,
			VisKm:      response.Current.VisKm,
			UV:         response.Current.UV,
		},
		Location: response.Location,
		Forecast: response.Forecast,
	}

	return weatherData, nil
}

// FetchHistoricalWeather retrieves historical weather data
func FetchHistoricalWeather(location string, date string) (*model.HistoricalData, error) {
	if config.AppConfig.APIKey == "" {
		return nil, fmt.Errorf("API key not set. Use --api-key flag or set ILLAPA_API_KEY environment variable")
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Fetching historical data "
	s.Start()
	defer s.Stop()

	// Construct historical URL
	historicalURL := fmt.Sprintf("%s/history.json?key=%s&q=%s&dt=%s",
		baseURL, config.AppConfig.APIKey, location, date)

	// Make request
	resp, err := http.Get(historicalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response weatherAPIHistoricalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Convert to our unified model
	historicalData := &model.HistoricalData{
		Location: response.Location,
		Forecast: response.Forecast,
	}

	return historicalData, nil
}

// WeatherAPI response models
type weatherAPIResponse struct {
	Location model.Location    `json:"location"`
	Current  weatherAPICurrent `json:"current"`
	Forecast model.Forecast    `json:"forecast"`
}

type weatherAPICurrent struct {
	LastUpdatedEpoch int64           `json:"last_updated_epoch"`
	LastUpdated      string          `json:"last_updated"`
	TempC            float64         `json:"temp_c"`
	TempF            float64         `json:"temp_f"`
	IsDay            int             `json:"is_day"`
	Condition        model.Condition `json:"condition"`
	WindMph          float64         `json:"wind_mph"`
	WindKph          float64         `json:"wind_kph"`
	WindDegree       int             `json:"wind_degree"`
	WindDir          string          `json:"wind_dir"`
	PressureMb       float64         `json:"pressure_mb"`
	PressureIn       float64         `json:"pressure_in"`
	PrecipMm         float64         `json:"precip_mm"`
	PrecipIn         float64         `json:"precip_in"`
	Humidity         int             `json:"humidity"`
	Cloud            int             `json:"cloud"`
	FeelslikeC       float64         `json:"feelslike_c"`
	FeelslikeF       float64         `json:"feelslike_f"`
	VisKm            float64         `json:"vis_km"`
	VisMiles         float64         `json:"vis_miles"`
	UV               float64         `json:"uv"`
	GustMph          float64         `json:"gust_mph"`
	GustKph          float64         `json:"gust_kph"`
}

type weatherAPIHistoricalResponse struct {
	Location model.Location `json:"location"`
	Forecast model.Forecast `json:"forecast"`
}
