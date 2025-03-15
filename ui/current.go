package ui

import (
	"fmt"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

// DisplayCurrentWeather outputs current weather conditions with a clean design
func DisplayCurrentWeather(data *model.WeatherData) {
	fmt.Println()

	// Location and current time with clean styling
	locationTitle := color.New(color.FgHiCyan, color.Bold)
	locationTitle.Printf("📍 %s, %s\n", data.Location.Name, data.Location.Country)
	fmt.Printf("🕒 Local time: %s\n", data.Location.Localtime)
	fmt.Println()

	// Current conditions with clean styling
	conditionIcon := GetConditionIcon(data.Current.Condition.Text)

	current := color.New(color.FgHiWhite, color.Bold)
	current.Println("Current Weather")
	fmt.Println()

	tempC := color.New(color.FgHiYellow, color.Bold)
	tempF := color.New(color.FgYellow)
	condition := color.New(color.FgHiWhite)

	condition.Printf("%s  %s ", conditionIcon, data.Current.Condition.Text)
	tempC.Printf("%.1f°C", data.Current.TempC)
	fmt.Printf(" / ")
	tempF.Printf("%.1f°F", data.Current.TempF)
	fmt.Println()

	feelsLike := color.New(color.FgHiWhite)
	feelsLike.Printf("Feels like: ")
	tempC.Printf("%.1f°C", data.Current.FeelsLikeC)
	fmt.Printf(" / ")
	tempF.Printf("%.1f°F", data.Current.FeelsLikeF)
	fmt.Println()
	fmt.Println()

	// Create styled labels for details
	labelStyle := color.New(color.FgHiBlue)
	valueStyle := color.New(color.FgWhite)

	// Wind info
	labelStyle.Printf("Wind:      ")
	valueStyle.Printf("%.1f km/h %s\n", data.Current.WindKph, data.Current.WindDir)

	// Humidity
	labelStyle.Printf("Humidity:  ")
	valueStyle.Printf("%d%%\n", data.Current.Humidity)

	// Precipitation
	labelStyle.Printf("Precip:    ")
	valueStyle.Printf("%.1f mm\n", data.Current.PrecipMm)

	// Visibility
	labelStyle.Printf("Visibility:")
	valueStyle.Printf(" %.1f km\n", data.Current.VisKm)

	// UV Index with color coding based on value
	labelStyle.Printf("UV Index:  ")

	// Color-code UV index based on intensity
	uvStyle := color.New(color.FgHiGreen)
	if data.Current.UV > 3 && data.Current.UV <= 6 {
		uvStyle = color.New(color.FgHiYellow)
	} else if data.Current.UV > 6 && data.Current.UV <= 8 {
		uvStyle = color.New(color.FgHiMagenta)
	} else if data.Current.UV > 8 {
		uvStyle = color.New(color.FgHiRed)
	}

	uvStyle.Printf("%.1f\n", data.Current.UV)
	fmt.Println()

	// Check alerts
	CheckAlerts(data)
}
