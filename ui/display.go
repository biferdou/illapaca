package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/biferdou/illapaca/config"
	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var weatherIcons = map[string]string{
	"sunny":         "☀️",
	"clear":         "☀️",
	"partly cloudy": "⛅",
	"cloudy":        "☁️",
	"overcast":      "☁️",
	"mist":          "🌫️",
	"rain":          "🌧️",
	"snow":          "❄️",
	"storm":         "⛈️",
	"fog":           "🌫️",
	// Add more mappings as needed
}

// getConditionIcon returns the appropriate weather icon
func getConditionIcon(condition string) string {
	conditionLower := strings.ToLower(condition)

	for k, v := range weatherIcons {
		if strings.Contains(conditionLower, k) {
			return v
		}
	}

	// Default icon if no match
	return "🌡️"
}

// DisplayCurrentWeather outputs current weather conditions
func DisplayCurrentWeather(data *model.WeatherData) {
	fmt.Println()

	// Location and current time
	locationTitle := color.New(color.FgHiCyan, color.Bold)
	locationTitle.Printf("📍 %s, %s, %s\n", data.Location.Name, data.Location.Region, data.Location.Country)
	fmt.Printf("🕒 Local time: %s\n\n", data.Location.Localtime)

	// Current conditions
	conditionIcon := getConditionIcon(data.Current.Condition.Text)

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

// DisplayForecast outputs weather forecast
func DisplayForecast(data *model.WeatherData) {
	forecastTitle := color.New(color.FgHiMagenta, color.Bold)
	forecastTitle.Println("Weather Forecast:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Condition", "Max", "Min", "Rain Chance", "Sunrise", "Sunset"})

	for _, day := range data.Forecast.ForecastDay {
		condition := day.Day.Condition.Text
		icon := getConditionIcon(condition)

		conditionWithIcon := fmt.Sprintf("%s %s", icon, condition)

		maxTemp := fmt.Sprintf("%.1f°C", day.Day.MaxTempC)
		minTemp := fmt.Sprintf("%.1f°C", day.Day.MinTempC)
		rainChance := fmt.Sprintf("%d%%", day.Day.DailyChanceOfRain)

		table.Append([]string{
			day.Date,
			conditionWithIcon,
			maxTemp,
			minTemp,
			rainChance,
			day.Astro.Sunrise,
			day.Astro.Sunset,
		})
	}

	table.Render()
	fmt.Println()
}

// DisplayTemperatureChart renders a simple ASCII chart of temperature trends
func DisplayTemperatureChart(data *model.WeatherData) {
	chartTitle := color.New(color.FgHiGreen, color.Bold)
	chartTitle.Println("Temperature Trend (24 hours):")

	// Get the first day's hourly forecast
	hours := data.Forecast.ForecastDay[0].Hour

	// Find max and min temps for scaling
	var maxTemp, minTemp float64
	for i, hour := range hours {
		if i == 0 || hour.TempC > maxTemp {
			maxTemp = hour.TempC
		}
		if i == 0 || hour.TempC < minTemp {
			minTemp = hour.TempC
		}
	}

	chartHeight := 10
	scale := float64(chartHeight) / (maxTemp - minTemp)

	// Print temperature scale
	fmt.Printf("%.1f°C │\n", maxTemp)
	fmt.Printf("%.1f°C │\n", minTemp)

	// Print chart
	for i := chartHeight; i >= 0; i-- {
		threshold := minTemp + float64(i)/scale

		fmt.Print("      │")
		for _, hour := range hours {
			if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
				if hour.TempC >= threshold {
					fmt.Print("●")
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}

	// Print time axis
	fmt.Print("      └")
	for _, hour := range hours {
		if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
			fmt.Print("─")
		}
	}
	fmt.Println()

	// Print hours
	fmt.Print("        ")
	for _, hour := range hours {
		if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
			t, _ := time.Parse("2006-01-02 15:04", hour.Time)
			fmt.Printf("%02d ", t.Hour())
		}
	}
	fmt.Println()
	fmt.Println()
}

// CheckAlerts checks weather against alert thresholds
func CheckAlerts(data *model.WeatherData) {
	alert := color.New(color.FgHiRed, color.Bold)
	thresholds := config.AppConfig.AlertThresholds

	// Check current conditions
	if data.Current.TempC > thresholds.HighTemp {
		alert.Printf("⚠️ ALERT: High temperature (%.1f°C) exceeds threshold (%.1f°C)\n",
			data.Current.TempC, thresholds.HighTemp)
	}

	if data.Current.TempC < thresholds.LowTemp {
		alert.Printf("⚠️ ALERT: Low temperature (%.1f°C) below threshold (%.1f°C)\n",
			data.Current.TempC, thresholds.LowTemp)
	}

	if data.Current.WindKph > thresholds.WindSpeed {
		alert.Printf("⚠️ ALERT: High wind speed (%.1f km/h) exceeds threshold (%.1f km/h)\n",
			data.Current.WindKph, thresholds.WindSpeed)
	}

	// Check forecast for precipitation
	for _, day := range data.Forecast.ForecastDay {
		if day.Day.DailyChanceOfRain > int(thresholds.Precipitation) {
			alert.Printf("⚠️ ALERT: High chance of rain (%d%%) on %s exceeds threshold (%.0f%%)\n",
				day.Day.DailyChanceOfRain, day.Date, thresholds.Precipitation)
		}
	}

	fmt.Println()
}

// DisplayHistoricalComparison shows comparison with historical data
func DisplayHistoricalComparison(current *model.WeatherData, historical *model.HistoricalData) {
	comparisonTitle := color.New(color.FgHiBlue, color.Bold)
	comparisonTitle.Printf("Historical Comparison (%s vs Today):\n",
		historical.Forecast.ForecastDay[0].Date)

	// Get current day's data
	currentDay := current.Forecast.ForecastDay[0].Day

	// Get historical day's data
	historicalDay := historical.Forecast.ForecastDay[0].Day

	// Display comparison
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Today", "Historical", "Difference"})

	// Max temperature
	maxTempDiff := currentDay.MaxTempC - historicalDay.MaxTempC
	maxTempDiffStr := fmt.Sprintf("%.1f°C", maxTempDiff)
	if maxTempDiff > 0 {
		maxTempDiffStr = "+" + maxTempDiffStr
	}

	// Min temperature
	minTempDiff := currentDay.MinTempC - historicalDay.MinTempC
	minTempDiffStr := fmt.Sprintf("%.1f°C", minTempDiff)
	if minTempDiff > 0 {
		minTempDiffStr = "+" + minTempDiffStr
	}

	// Precipitation
	precipDiff := currentDay.TotalPrecipMm - historicalDay.TotalPrecipMm
	precipDiffStr := fmt.Sprintf("%.1f mm", precipDiff)
	if precipDiff > 0 {
		precipDiffStr = "+" + precipDiffStr
	}

	table.Append([]string{"Max Temp",
		fmt.Sprintf("%.1f°C", currentDay.MaxTempC),
		fmt.Sprintf("%.1f°C", historicalDay.MaxTempC),
		maxTempDiffStr})

	table.Append([]string{"Min Temp",
		fmt.Sprintf("%.1f°C", currentDay.MinTempC),
		fmt.Sprintf("%.1f°C", historicalDay.MinTempC),
		minTempDiffStr})

	table.Append([]string{"Precipitation",
		fmt.Sprintf("%.1f mm", currentDay.TotalPrecipMm),
		fmt.Sprintf("%.1f mm", historicalDay.TotalPrecipMm),
		precipDiffStr})

	table.Render()
	fmt.Println()
}

// DisplayDashboard displays the full dashboard
func DisplayDashboard(data *model.WeatherData) {
	// Title banner
	title := color.New(color.FgHiCyan, color.Bold)
	title.Println("╔══════════════════════════════════════════╗")
	title.Println("║           ILLAPA WEATHER DASHBOARD       ║")
	title.Println("╚══════════════════════════════════════════╝")

	DisplayCurrentWeather(data)
	DisplayForecast(data)
	DisplayTemperatureChart(data)
}
