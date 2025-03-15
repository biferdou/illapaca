package ui

import (
	"fmt"

	"github.com/biferdou/illapaca/config"
	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

// DisplayDashboard displays the full dashboard
func DisplayDashboard(data *model.WeatherData) {
	// Title banner
	displayDashboardHeader()

	// Display components in sequence
	DisplayCurrentWeather(data)
	DisplayForecast(data)
	DisplayTemperatureChart(data)

	// If at least one day of forecast available, show precipitation chart
	if len(data.Forecast.ForecastDay) > 0 {
		DisplayPrecipitationChart(data.Forecast.ForecastDay[0])
	}
}

// displayDashboardHeader displays the dashboard title banner
func displayDashboardHeader() {
	title := color.New(color.FgHiCyan, color.Bold)
	title.Println("╔══════════════════════════════════════════╗")
	title.Println("║           ILLAPA WEATHER DASHBOARD       ║")
	title.Println("╚══════════════════════════════════════════╝")
}

// DisplayExtendedDashboard shows a more detailed dashboard with hourly forecasts
func DisplayExtendedDashboard(data *model.WeatherData, showHourly bool) {
	displayDashboardHeader()

	DisplayCurrentWeather(data)
	DisplayForecast(data)
	DisplayTemperatureChart(data)

	if len(data.Forecast.ForecastDay) > 0 {
		// Show today's precipitation chart
		DisplayPrecipitationChart(data.Forecast.ForecastDay[0])

		// Show hourly forecast if requested
		if showHourly {
			DisplayHourlyForecast(data.Forecast.ForecastDay[0])
		}
	}
}

// DisplayCompactDashboard shows a minimal dashboard for small terminals
func DisplayCompactDashboard(data *model.WeatherData) {
	compactTitle := color.New(color.FgHiCyan, color.Bold)
	compactTitle.Println("=== ILLAPA WEATHER ===")

	// Simplified current weather display
	locationTitle := color.New(color.FgHiCyan)
	locationTitle.Printf("📍 %s, %s | %s\n",
		data.Location.Name, data.Location.Country, data.Location.Localtime)

	// Current conditions - compact format
	conditionIcon := GetConditionIcon(data.Current.Condition.Text)

	tempC := color.New(color.FgHiYellow, color.Bold)
	fmt.Printf("%s %s ", conditionIcon, data.Current.Condition.Text)
	tempC.Printf("%.1f°C", data.Current.TempC)
	fmt.Printf(" (Feels: %.1f°C) | ", data.Current.FeelsLikeC)
	fmt.Printf("Wind: %.1f km/h %s | Hum: %d%%\n\n",
		data.Current.WindKph, data.Current.WindDir, data.Current.Humidity)

	// Compact forecast
	forecastTitle := color.New(color.FgHiMagenta)
	forecastTitle.Println("3-Day Forecast:")

	days := min(len(data.Forecast.ForecastDay), 3)

	for i := range days {
		day := data.Forecast.ForecastDay[i]
		icon := GetConditionIcon(day.Day.Condition.Text)
		fmt.Printf("%s: %s %.1f°C/%.1f°C | Rain: %d%%\n",
			day.Date, icon, day.Day.MaxTempC, day.Day.MinTempC, day.Day.DailyChanceOfRain)
	}

	fmt.Println()

	// Check alerts but only show count
	alerts := getAlerts(data, config.AppConfig.AlertThresholds)
	if len(alerts) > 0 {
		alertMsg := color.New(color.FgHiRed, color.Bold)
		alertMsg.Printf("⚠️ %d weather alerts detected\n\n", len(alerts))
	}
}
