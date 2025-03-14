package ui

import (
	"fmt"
	"os"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// DisplayForecast outputs weather forecast
func DisplayForecast(data *model.WeatherData) {
	forecastTitle := color.New(color.FgHiMagenta, color.Bold)
	forecastTitle.Println("Weather Forecast:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Condition", "Max", "Min", "Rain Chance", "Sunrise", "Sunset"})

	for _, day := range data.Forecast.ForecastDay {
		condition := day.Day.Condition.Text
		icon := GetConditionIcon(condition)

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

// DisplayHourlyForecast outputs hourly weather forecast for a given day
func DisplayHourlyForecast(day model.ForecastDay) {
	hourlyTitle := color.New(color.FgHiCyan, color.Bold)
	hourlyTitle.Printf("Hourly Forecast for %s:\n", day.Date)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Temp", "Condition", "Rain Chance"})

	// Display only a subset of hours to keep the output manageable
	for i, hour := range day.Hour {
		if i%3 != 0 { // Skip to show only every 3 hours
			continue
		}

		// Extract just the time portion (15:04)
		timeOnly := hour.Time[11:16]

		temp := fmt.Sprintf("%.1f°C", hour.TempC)
		condition := GetConditionIcon(hour.Condition.Text) + " " + hour.Condition.Text
		rainChance := fmt.Sprintf("%d%%", hour.ChanceOfRain)

		table.Append([]string{
			timeOnly,
			temp,
			condition,
			rainChance,
		})
	}

	table.Render()
	fmt.Println()
}
