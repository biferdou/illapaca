package ui

import (
	"fmt"

	"github.com/biferdou/illapaca/config"
	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

// CheckAlerts checks weather against alert thresholds with improved styling
func CheckAlerts(data *model.WeatherData) {
	thresholds := config.AppConfig.AlertThresholds
	alerts := getAlerts(data, thresholds)

	if len(alerts) == 0 {
		return
	}

	// Create an alert box
	alertBox := color.New(color.FgHiWhite)
	alertTitle := color.New(color.FgHiRed, color.Bold)
	alertText := color.New(color.FgHiRed)

	alertBox.Println("┌─────────────────────────────────────────┐")
	alertBox.Print("│ ")
	alertTitle.Println("⚠️  WEATHER ALERTS  ⚠️")

	// Display each alert
	for _, alertMsg := range alerts {
		alertBox.Print("│ ")
		alertText.Printf("• %s", alertMsg)
		alertBox.Println(" │")
	}

	alertBox.Println("└─────────────────────────────────────────┘")
	fmt.Println()
}

// getAlerts returns a list of alert messages for the given weather data
func getAlerts(data *model.WeatherData, thresholds config.AlertThresholds) []string {
	var alerts []string

	// Check current conditions
	if data.Current.TempC > thresholds.HighTemp {
		alerts = append(alerts, fmt.Sprintf("High temperature (%.1f°C) exceeds threshold (%.1f°C)",
			data.Current.TempC, thresholds.HighTemp))
	}

	if data.Current.TempC < thresholds.LowTemp {
		alerts = append(alerts, fmt.Sprintf("Low temperature (%.1f°C) below threshold (%.1f°C)",
			data.Current.TempC, thresholds.LowTemp))
	}

	if data.Current.WindKph > thresholds.WindSpeed {
		alerts = append(alerts, fmt.Sprintf("High wind speed (%.1f km/h) exceeds threshold (%.1f km/h)",
			data.Current.WindKph, thresholds.WindSpeed))
	}

	// Check forecast for precipitation
	for _, day := range data.Forecast.ForecastDay {
		if day.Day.DailyChanceOfRain > int(thresholds.Precipitation) {
			alerts = append(alerts, fmt.Sprintf("High chance of rain (%d%%) on %s exceeds threshold (%.0f%%)",
				day.Day.DailyChanceOfRain, day.Date, thresholds.Precipitation))
		}
	}

	return alerts
}

// DisplayAlertSettings shows the current alert threshold settings
func DisplayAlertSettings(thresholds config.AlertThresholds) {
	settingsTitle := color.New(color.FgHiYellow, color.Bold)
	settingsTitle.Println("Alert Threshold Settings:")

	fmt.Printf("High Temperature: %.1f°C\n", thresholds.HighTemp)
	fmt.Printf("Low Temperature: %.1f°C\n", thresholds.LowTemp)
	fmt.Printf("Precipitation Chance: %.0f%%\n", thresholds.Precipitation)
	fmt.Printf("Wind Speed: %.1f km/h\n", thresholds.WindSpeed)
	fmt.Println()
}
