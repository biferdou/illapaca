package ui

import (
	"fmt"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

// DisplayCurrentWeather outputs current weather conditions
func DisplayCurrentWeather(data *model.WeatherData) {
	fmt.Println()

	// Location and current time
	locationTitle := color.New(color.FgHiCyan, color.Bold)
	locationTitle.Printf("📍 %s, %s\n", data.Location.Name, data.Location.Country)
	fmt.Printf("🕒 Local time: %s\n\n", data.Location.Localtime)

	// Current conditions
	conditionIcon := GetConditionIcon(data.Current.Condition.Text)

	current := color.New(color.FgHiWhite, color.Bold)
	current.Println("Current Weather:")

	tempC := color.New(color.FgHiYellow, color.Bold)
	tempF := color.New(color.FgYellow)

	fmt.Printf("%s %s ", conditionIcon, data.Current.Condition.Text)
	tempC.Printf("%.1f°C", data.Current.TempC)
	fmt.Printf(" / ")
	tempF.Printf("%.1f°F", data.Current.TempF)
	fmt.Println()

	fmt.Printf("Feels like: %.1f°C / %.1f°F\n", data.Current.FeelsLikeC, data.Current.FeelsLikeF)
	fmt.Printf("Wind: %.1f km/h %s\n", data.Current.WindKph, data.Current.WindDir)
	fmt.Printf("Humidity: %d%%\n", data.Current.Humidity)
	fmt.Printf("Precipitation: %.1f mm\n", data.Current.PrecipMm)
	fmt.Printf("Visibility: %.1f km\n", data.Current.VisKm)
	fmt.Printf("UV Index: %.1f\n", data.Current.UV)

	fmt.Println()

	// Check alerts
	CheckAlerts(data)
}
