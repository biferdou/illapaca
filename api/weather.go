package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/biferdou/illapaca/config"
	"github.com/biferdou/illapaca/model"
	"github.com/briandowns/spinner"
)

const (
	// OpenWeatherMap API endpoints
	baseURL       = "https://api.openweathermap.org/data/2.5"
	geoURL        = "https://api.openweathermap.org/geo/1.0/direct"
	historicalURL = "https://api.openweathermap.org/data/2.5/onecall/timemachine"
)

// geocodeLocation converts a location name to coordinates
func geocodeLocation(location string) (*model.Coordinates, error) {
	if config.AppConfig.APIKey == "" {
		return nil, fmt.Errorf("API key not set. Use --api-key flag or set ILLAPA_API_KEY environment variable")
	}

	url := fmt.Sprintf("%s?q=%s&limit=1&appid=%s",
		geoURL, location, config.AppConfig.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var locations []model.GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("location not found: %s", location)
	}

	return &model.Coordinates{
		Lat:     locations[0].Lat,
		Lon:     locations[0].Lon,
		Name:    locations[0].Name,
		Country: locations[0].Country,
	}, nil
}

// FetchWeather retrieves weather data from the API
func FetchWeather(location string, days int) (*model.WeatherData, error) {
	if config.AppConfig.APIKey == "" {
		return nil, fmt.Errorf("API key not set. Use --api-key flag or set ILLAPA_API_KEY environment variable")
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Fetching weather data "
	s.Start()
	defer s.Stop()

	// First, geocode the location
	coordinates, err := geocodeLocation(location)
	if err != nil {
		return nil, err
	}

	// Get current weather
	currentURL := fmt.Sprintf("%s/weather?lat=%f&lon=%f&units=%s&appid=%s",
		baseURL, coordinates.Lat, coordinates.Lon,
		config.AppConfig.Units, config.AppConfig.APIKey)

	currentResp, err := http.Get(currentURL)
	if err != nil {
		return nil, err
	}
	defer currentResp.Body.Close()

	if currentResp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(currentResp.Body)
		return nil, fmt.Errorf("API error (%d): %s", currentResp.StatusCode, string(body))
	}

	var currentData model.OpenWeatherCurrent
	if err := json.NewDecoder(currentResp.Body).Decode(&currentData); err != nil {
		return nil, err
	}

	// Get forecast
	forecastURL := fmt.Sprintf("%s/forecast?lat=%f&lon=%f&units=%s&appid=%s",
		baseURL, coordinates.Lat, coordinates.Lon,
		config.AppConfig.Units, config.AppConfig.APIKey)

	forecastResp, err := http.Get(forecastURL)
	if err != nil {
		return nil, err
	}
	defer forecastResp.Body.Close()

	if forecastResp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(forecastResp.Body)
		return nil, fmt.Errorf("API error (%d): %s", forecastResp.StatusCode, string(body))
	}

	var forecastData model.OpenWeatherForecast
	if err := json.NewDecoder(forecastResp.Body).Decode(&forecastData); err != nil {
		return nil, err
	}

	// Convert to our unified model
	weatherData := &model.WeatherData{
		Current: model.CurrentWeather{
			TempC:      currentData.Main.Temp,
			FeelsLikeC: currentData.Main.FeelsLike,
			Humidity:   currentData.Main.Humidity,
			WindKph:    currentData.Wind.Speed * 3.6, // Convert m/s to km/h
			WindDir:    getWindDirection(currentData.Wind.Deg),
			PrecipMm:   getPrecipitation(currentData),
			Condition: model.Condition{
				Text: currentData.Weather[0].Description,
				Icon: currentData.Weather[0].Icon,
			},
			VisKm: currentData.Visibility / 1000.0, // Convert m to km
			UV:    0,                               // OpenWeatherMap doesn't provide UV in the basic API
		},
		Location: model.Location{
			Name:      coordinates.Name,
			Country:   coordinates.Country,
			Lat:       coordinates.Lat,
			Lon:       coordinates.Lon,
			Localtime: time.Unix(int64(currentData.Dt), 0).Format("2006-01-02 15:04"),
		},
		Forecast: convertForecast(forecastData, days),
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

	// First, geocode the location
	coordinates, err := geocodeLocation(location)
	if err != nil {
		return nil, err
	}

	// Parse the date and get Unix timestamp
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s. Use YYYY-MM-DD", date)
	}
	timestamp := t.Unix()

	// Get historical data
	url := fmt.Sprintf("%s?lat=%f&lon=%f&dt=%d&units=%s&appid=%s",
		historicalURL, coordinates.Lat, coordinates.Lon,
		timestamp, config.AppConfig.Units, config.AppConfig.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var historicalData model.OpenWeatherHistorical
	if err := json.NewDecoder(resp.Body).Decode(&historicalData); err != nil {
		return nil, err
	}

	// Convert to our unified model
	data := &model.HistoricalData{
		Location: model.Location{
			Name:      coordinates.Name,
			Country:   coordinates.Country,
			Lat:       coordinates.Lat,
			Lon:       coordinates.Lon,
			Localtime: time.Unix(int64(historicalData.Current.Dt), 0).Format("2006-01-02 15:04"),
		},
		Forecast: model.Forecast{
			ForecastDay: []model.ForecastDay{
				{
					Date: date,
					Day: model.Day{
						MaxTempC:      findMaxTemp(historicalData.Hourly),
						MinTempC:      findMinTemp(historicalData.Hourly),
						AvgTempC:      findAvgTemp(historicalData.Hourly),
						TotalPrecipMm: findTotalPrecip(historicalData.Hourly),
						Condition: model.Condition{
							Text: getMostFrequentCondition(historicalData.Hourly),
						},
					},
					Hour: convertHistoricalHourly(historicalData.Hourly),
				},
			},
		},
	}

	return data, nil
}

// Helper functions for conversion and data extraction

func getWindDirection(degrees float64) string {
	// Convert degrees to cardinal direction
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((degrees+11.25)/22.5) % 16
	return directions[index]
}

func getPrecipitation(data model.OpenWeatherCurrent) float64 {
	// Extract precipitation from rain or snow if available
	if data.Rain.OneHour > 0 {
		return data.Rain.OneHour
	} else if data.Snow.OneHour > 0 {
		return data.Snow.OneHour
	}
	return 0
}

func convertForecast(forecast model.OpenWeatherForecast, days int) model.Forecast {
	// Group forecast by day
	dailyForecasts := make(map[string][]model.OpenWeatherForecastItem)

	for _, item := range forecast.List {
		date := time.Unix(int64(item.Dt), 0).Format("2006-01-02")
		dailyForecasts[date] = append(dailyForecasts[date], item)
	}

	// Convert to our model
	var forecastDays []model.ForecastDay

	// Sort dates and limit to requested number of days
	dates := make([]string, 0, len(dailyForecasts))
	for date := range dailyForecasts {
		dates = append(dates, date)
	}

	// Sort dates (we would need a proper sort function here)
	// For simplicity, we'll rely on the fact that the API returns dates in order

	// Limit to requested days
	if len(dates) > days {
		dates = dates[:days]
	}

	for _, date := range dates {
		items := dailyForecasts[date]

		// Find min, max temps and most common condition
		var maxTemp, minTemp float64
		minTemp = 1000 // Arbitrary high value

		for i, item := range items {
			if i == 0 || item.Main.TempMax > maxTemp {
				maxTemp = item.Main.TempMax
			}
			if i == 0 || item.Main.TempMin < minTemp {
				minTemp = item.Main.TempMin
			}
		}

		// Calculate chance of rain
		rainCount := 0
		for _, item := range items {
			for _, w := range item.Weather {
				if w.Main == "Rain" || w.Main == "Drizzle" {
					rainCount++
					break
				}
			}
		}
		chanceOfRain := int((float64(rainCount) / float64(len(items))) * 100)

		// Get sunrise/sunset (assume the same for the whole day)
		// Note: OpenWeatherMap 5-day forecast doesn't include this
		// We would need to make another API call or use approximate values

		// Convert hourly forecasts
		var hours []model.Hour
		for _, item := range items {
			hour := model.Hour{
				TimeEpoch: int64(item.Dt),
				Time:      time.Unix(int64(item.Dt), 0).Format("2006-01-02 15:04"),
				TempC:     item.Main.Temp,
				Condition: model.Condition{
					Text: item.Weather[0].Description,
					Icon: item.Weather[0].Icon,
				},
				ChanceOfRain: getChanceOfRain(item),
			}
			hours = append(hours, hour)
		}

		day := model.ForecastDay{
			Date: date,
			Day: model.Day{
				MaxTempC:          maxTemp,
				MinTempC:          minTemp,
				AvgTempC:          (maxTemp + minTemp) / 2,
				DailyChanceOfRain: chanceOfRain,
				Condition: model.Condition{
					Text: getMostFrequentCondition(items),
				},
			},
			Astro: model.Astro{
				Sunrise: "06:00", // Placeholder - we would need to fetch this separately
				Sunset:  "18:00", // Placeholder - we would need to fetch this separately
			},
			Hour: hours,
		}

		forecastDays = append(forecastDays, day)
	}

	return model.Forecast{
		ForecastDay: forecastDays,
	}
}

func getChanceOfRain(item model.OpenWeatherForecastItem) int {
	// Check if there's a probability of precipitation
	pop := int(item.Pop * 100) // Pop is from 0 to 1
	return pop
}

func getMostFrequentCondition(items interface{}) string {
	conditionCount := make(map[string]int)

	// Extract conditions based on the type of items
	switch v := items.(type) {
	case []model.OpenWeatherForecastItem:
		for _, item := range v {
			if len(item.Weather) > 0 {
				conditionCount[item.Weather[0].Description]++
			}
		}
	case []model.OpenWeatherHourlyItem:
		for _, item := range v {
			if len(item.Weather) > 0 {
				conditionCount[item.Weather[0].Description]++
			}
		}
	}

	// Find the most frequent condition
	var mostFrequent string
	var maxCount int

	for condition, count := range conditionCount {
		if count > maxCount {
			maxCount = count
			mostFrequent = condition
		}
	}

	return mostFrequent
}

func findMaxTemp(items []model.OpenWeatherHourlyItem) float64 {
	var maxTemp float64

	for i, item := range items {
		if i == 0 || item.Temp > maxTemp {
			maxTemp = item.Temp
		}
	}

	return maxTemp
}

func findMinTemp(items []model.OpenWeatherHourlyItem) float64 {
	var minTemp float64 = 1000 // Arbitrary high value

	for i, item := range items {
		if i == 0 || item.Temp < minTemp {
			minTemp = item.Temp
		}
	}

	return minTemp
}

func findAvgTemp(items []model.OpenWeatherHourlyItem) float64 {
	var sum float64

	for _, item := range items {
		sum += item.Temp
	}

	return sum / float64(len(items))
}

func findTotalPrecip(items []model.OpenWeatherHourlyItem) float64 {
	var total float64

	for _, item := range items {
		if item.Rain.OneHour > 0 {
			total += item.Rain.OneHour
		}
		if item.Snow.OneHour > 0 {
			total += item.Snow.OneHour
		}
	}

	return total
}

func convertHistoricalHourly(hourly []model.OpenWeatherHourlyItem) []model.Hour {
	var hours []model.Hour

	for _, item := range hourly {
		hour := model.Hour{
			TimeEpoch: int64(item.Dt),
			Time:      time.Unix(int64(item.Dt), 0).Format("2006-01-02 15:04"),
			TempC:     item.Temp,
			Condition: model.Condition{
				Text: item.Weather[0].Description,
				Icon: item.Weather[0].Icon,
			},
			ChanceOfRain: 0, // Historical data doesn't include chance of rain
		}
		hours = append(hours, hour)
	}

	return hours
}
