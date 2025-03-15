package ui

import (
	"fmt"
	"os"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// DisplayForecast outputs weather forecast with clean styling
func DisplayForecast(data *model.WeatherData) {
	forecastTitle := color.New(color.FgHiMagenta, color.Bold)
	forecastTitle.Println("Weather Forecast")
	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Condition", "Max", "Min", "Rain", "Sunrise", "Sunset"})
	// Ensure the table has a consistent width by setting column alignments
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("  ") // Use double spaces instead of pipes
	table.SetRowSeparator("")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor},
	)

	for _, day := range data.Forecast.ForecastDay {
		condition := day.Day.Condition.Text
		icon := GetConditionIcon(condition)

		// Use just the icon for the condition to save space and maintain alignment
		conditionWithIcon := icon

		maxTemp := fmt.Sprintf("%.1f°C", day.Day.MaxTempC)
		minTemp := fmt.Sprintf("%.1f°C", day.Day.MinTempC)

		// Style rain chance based on probability
		var rainChance string
		rainProb := day.Day.DailyChanceOfRain
		if rainProb < 20 {
			rainChance = fmt.Sprintf("\x1b[38;5;39m%d%%\x1b[0m", rainProb)
		} else if rainProb < 40 {
			rainChance = fmt.Sprintf("\x1b[38;5;45m%d%%\x1b[0m", rainProb)
		} else if rainProb < 60 {
			rainChance = fmt.Sprintf("\x1b[38;5;51m%d%%\x1b[0m", rainProb)
		} else if rainProb < 80 {
			rainChance = fmt.Sprintf("\x1b[38;5;33m%d%%\x1b[0m", rainProb)
		} else {
			rainChance = fmt.Sprintf("\x1b[38;5;27m%d%%\x1b[0m", rainProb)
		}

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
	hourlyTitle.Printf("Hourly Forecast for %s\n", day.Date)
	fmt.Println()

	// Create a lookup table for condition descriptions
	conditionDescriptions := make(map[string]string)
	for i, hour := range day.Hour {
		if i%3 == 0 { // Only include every 3 hours to save space
			icon := GetConditionIcon(hour.Condition.Text)
			conditionDescriptions[icon] = hour.Condition.Text
		}
	}

	// Display condition key first
	fmt.Println("Weather conditions:")
	for icon, description := range conditionDescriptions {
		fmt.Printf("%s %s\n", icon, description)
	}
	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Temp", "Condition", "Rain Chance"})
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("  ") // Use double spaces instead of pipes
	table.SetRowSeparator("")

	// Display only a subset of hours to keep the output manageable
	for i, hour := range day.Hour {
		if i%3 != 0 { // Skip to show only every 3 hours
			continue
		}

		// Extract just the time portion (15:04)
		timeOnly := hour.Time[11:16]

		temp := fmt.Sprintf("%.1f°C", hour.TempC)
		// Use just the icon for display, not the full condition text
		condition := GetConditionIcon(hour.Condition.Text)
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
