package ui

import (
	"fmt"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

// DisplayCurrentWeather outputs current weather conditions
func DisplayCurrentWeather(data *model.WeatherData) {
	fmt.Println()

	// Location and current time with improved styling
	locationBox := color.New(color.FgHiWhite)
	locationTitle := color.New(color.FgHiCyan, color.Bold)

	locationBox.Println("┌─────────────────────────────────────────┐")
	locationBox.Print("│ ")
	locationTitle.Printf("📍 %s, %s", data.Location.Name, data.Location.Country)
	locationBox.Println(" │")

	locationBox.Print("│ ")
	fmt.Printf("🕒 Local time: %s", data.Location.Localtime)
	locationBox.Println(" │")
	locationBox.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	// Current conditions with improved styling
	conditionIcon := GetConditionIcon(data.Current.Condition.Text)

	current := color.New(color.FgHiWhite, color.Bold)
	currentBox := color.New(color.FgHiWhite)

	currentBox.Println("┌─────────────────────────────────────────┐")
	currentBox.Print("│ ")
	current.Println("Current Weather:")
	currentBox.Print("│ ")

	tempC := color.New(color.FgHiYellow, color.Bold)
	tempF := color.New(color.FgYellow)
	condition := color.New(color.FgHiWhite)

	condition.Printf("%s %s ", conditionIcon, data.Current.Condition.Text)
	tempC.Printf("%.1f°C", data.Current.TempC)
	fmt.Printf(" / ")
	tempF.Printf("%.1f°F", data.Current.TempF)
	currentBox.Println(" │")

	currentBox.Print("│ ")
	feelsLike := color.New(color.FgHiWhite)
	feelsLike.Printf("Feels like: ")
	tempC.Printf("%.1f°C", data.Current.FeelsLikeC)
	fmt.Printf(" / ")
	tempF.Printf("%.1f°F", data.Current.FeelsLikeF)
	currentBox.Println(" │")

	// Create styled labels for details
	labelStyle := color.New(color.FgHiBlue)
	valueStyle := color.New(color.FgWhite)

	// Wind info
	currentBox.Print("│ ")
	labelStyle.Printf("Wind:      ")
	valueStyle.Printf("%.1f km/h %s", data.Current.WindKph, data.Current.WindDir)
	currentBox.Println(" │")

	// Humidity
	currentBox.Print("│ ")
	labelStyle.Printf("Humidity:  ")
	valueStyle.Printf("%d%%", data.Current.Humidity)
	currentBox.Println(" │")

	// Precipitation
	currentBox.Print("│ ")
	labelStyle.Printf("Precip:    ")
	valueStyle.Printf("%.1f mm", data.Current.PrecipMm)
	currentBox.Println(" │")

	// Visibility
	currentBox.Print("│ ")
	labelStyle.Printf("Visibility:")
	valueStyle.Printf(" %.1f km", data.Current.VisKm)
	currentBox.Println(" │")

	// UV Index with color coding based on value
	currentBox.Print("│ ")
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

	uvStyle.Printf("%.1f", data.Current.UV)
	currentBox.Println(" │")

	currentBox.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	// Check alerts
	CheckAlerts(data)
}
