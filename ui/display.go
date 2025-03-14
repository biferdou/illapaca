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
	locationTitle.Printf("📍 %s, %s\n", data.Location.Name, data.Location.Country)
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

// DisplayLocationComparison shows a side-by-side comparison of two locations
func DisplayLocationComparison(data1, data2 *model.WeatherData) {
	comparisonTitle := color.New(color.FgHiBlue, color.Bold)
	comparisonTitle.Printf("Location Comparison: %s vs %s\n\n",
		data1.Location.Name, data2.Location.Name)

	// Display comparison
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", data1.Location.Name, data2.Location.Name, "Difference"})

	// Temperature
	tempDiff := data1.Current.TempC - data2.Current.TempC
	tempDiffStr := fmt.Sprintf("%.1f°C", tempDiff)
	if tempDiff > 0 {
		tempDiffStr = "+" + tempDiffStr
	}

	// Feels like
	feelsLikeDiff := data1.Current.FeelsLikeC - data2.Current.FeelsLikeC
	feelsLikeDiffStr := fmt.Sprintf("%.1f°C", feelsLikeDiff)
	if feelsLikeDiff > 0 {
		feelsLikeDiffStr = "+" + feelsLikeDiffStr
	}

	// Humidity
	humidityDiff := data1.Current.Humidity - data2.Current.Humidity
	humidityDiffStr := fmt.Sprintf("%d%%", humidityDiff)
	if humidityDiff > 0 {
		humidityDiffStr = "+" + humidityDiffStr
	}

	// Wind speed
	windDiff := data1.Current.WindKph - data2.Current.WindKph
	windDiffStr := fmt.Sprintf("%.1f km/h", windDiff)
	if windDiff > 0 {
		windDiffStr = "+" + windDiffStr
	}

	// Build table rows
	table.Append([]string{"Temperature",
		fmt.Sprintf("%.1f°C", data1.Current.TempC),
		fmt.Sprintf("%.1f°C", data2.Current.TempC),
		tempDiffStr})

	table.Append([]string{"Feels Like",
		fmt.Sprintf("%.1f°C", data1.Current.FeelsLikeC),
		fmt.Sprintf("%.1f°C", data2.Current.FeelsLikeC),
		feelsLikeDiffStr})

	table.Append([]string{"Condition",
		data1.Current.Condition.Text,
		data2.Current.Condition.Text,
		"--"})

	table.Append([]string{"Humidity",
		fmt.Sprintf("%d%%", data1.Current.Humidity),
		fmt.Sprintf("%d%%", data2.Current.Humidity),
		humidityDiffStr})

	table.Append([]string{"Wind Speed",
		fmt.Sprintf("%.1f km/h", data1.Current.WindKph),
		fmt.Sprintf("%.1f km/h", data2.Current.WindKph),
		windDiffStr})

	table.Append([]string{"Wind Direction",
		data1.Current.WindDir,
		data2.Current.WindDir,
		"--"})

	table.Append([]string{"Precipitation",
		fmt.Sprintf("%.1f mm", data1.Current.PrecipMm),
		fmt.Sprintf("%.1f mm", data2.Current.PrecipMm),
		"--"})

	table.Append([]string{"Visibility",
		fmt.Sprintf("%.1f km", data1.Current.VisKm),
		fmt.Sprintf("%.1f km", data2.Current.VisKm),
		"--"})

	// Time difference
	table.Append([]string{"Local Time",
		data1.Location.Localtime,
		data2.Location.Localtime,
		"--"})

	table.Render()
	fmt.Println()

	// Add some analysis
	if tempDiff > 3 {
		fmt.Printf("📊 %s is %.1f°C warmer than %s\n",
			data1.Location.Name, tempDiff, data2.Location.Name)
	} else if tempDiff < -3 {
		fmt.Printf("📊 %s is %.1f°C colder than %s\n",
			data1.Location.Name, -tempDiff, data2.Location.Name)
	}

	if humidityDiff > 15 {
		fmt.Printf("💧 %s is more humid than %s\n",
			data1.Location.Name, data2.Location.Name)
	} else if humidityDiff < -15 {
		fmt.Printf("💧 %s is drier than %s\n",
			data1.Location.Name, data2.Location.Name)
	}

	if windDiff > 10 {
		fmt.Printf("🌬️ %s is windier than %s\n",
			data1.Location.Name, data2.Location.Name)
	} else if windDiff < -10 {
		fmt.Printf("🌬️ %s is calmer than %s\n",
			data1.Location.Name, data2.Location.Name)
	}
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
