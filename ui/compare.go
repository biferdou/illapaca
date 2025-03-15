package ui

import (
	"fmt"
	"os"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// DisplayLocationComparison shows a side-by-side comparison of two locations
func DisplayLocationComparison(data1, data2 *model.WeatherData) {
	// Create styled header
	compareBox := color.New(color.FgHiWhite)
	comparisonTitle := color.New(color.FgHiBlue, color.Bold)
	locationStyle := color.New(color.FgHiCyan, color.Bold)

	compareBox.Println("┌─────────────────────────────────────────────────────┐")
	compareBox.Print("│ ")
	comparisonTitle.Print("Location Comparison: ")
	locationStyle.Printf("%s", data1.Location.Name)
	comparisonTitle.Print(" vs ")
	locationStyle.Printf("%s", data2.Location.Name)
	compareBox.Println(" │")
	compareBox.Println("└─────────────────────────────────────────────────────┘")
	fmt.Println()

	// Display icon and current conditions in a header row
	displayQuickComparison(data1, data2)

	// Display comparison table
	table := createComparisonTable(data1, data2)
	table.Render()
	fmt.Println()

	// Add some analysis
	displayComparisonAnalysis(data1, data2)
}

// displayQuickComparison shows a quick visual comparison of current conditions
func displayQuickComparison(data1, data2 *model.WeatherData) {
	// Create styled boxes for each location
	boxStyle := color.New(color.FgHiWhite)
	tempStyle := color.New(color.FgHiYellow, color.Bold)

	// Get weather icons
	icon1 := GetConditionIcon(data1.Current.Condition.Text)
	icon2 := GetConditionIcon(data2.Current.Condition.Text)

	// Create weather condition lookup
	fmt.Println("Weather conditions:")
	fmt.Printf("%s %s\n", icon1, data1.Current.Condition.Text)
	fmt.Printf("%s %s\n", icon2, data2.Current.Condition.Text)
	fmt.Println()

	// Get condition text, limit to 22 characters to prevent display issues
	cond1 := data1.Current.Condition.Text
	if len(cond1) > 22 {
		cond1 = cond1[:19] + "..."
	}

	cond2 := data2.Current.Condition.Text
	if len(cond2) > 22 {
		cond2 = cond2[:19] + "..."
	}

	// Location 1 quick view
	boxStyle.Println("┌─────────────────────────────┐ ┌─────────────────────────────┐")
	boxStyle.Printf("│ %-27s │ │ %-27s │\n", data1.Location.Name, data2.Location.Name)
	boxStyle.Printf("│ %-2s %-24s │ │ %-2s %-24s │\n",
		icon1, cond1, icon2, cond2)

	boxStyle.Print("│ ")
	tempStyle.Printf("%.1f°C", data1.Current.TempC)
	boxStyle.Printf(" feels like ")
	tempStyle.Printf("%.1f°C", data1.Current.FeelsLikeC)

	boxStyle.Print(" │ │ ")
	tempStyle.Printf("%.1f°C", data2.Current.TempC)
	boxStyle.Printf(" feels like ")
	tempStyle.Printf("%.1f°C", data2.Current.FeelsLikeC)
	boxStyle.Println(" │")

	boxStyle.Println("└─────────────────────────────┘ └─────────────────────────────┘")
	fmt.Println()
}

// createComparisonTable builds the comparison table for two locations
func createComparisonTable(data1, data2 *model.WeatherData) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", data1.Location.Name, data2.Location.Name, "Difference"})
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor},
	)

	// Add a row for condition with just the icons
	table.Append([]string{
		"Condition",
		GetConditionIcon(data1.Current.Condition.Text),
		GetConditionIcon(data2.Current.Condition.Text),
		"--"})

	// Calculate differences
	tempDiff := data1.Current.TempC - data2.Current.TempC
	feelsLikeDiff := data1.Current.FeelsLikeC - data2.Current.FeelsLikeC
	humidityDiff := data1.Current.Humidity - data2.Current.Humidity
	windDiff := data1.Current.WindKph - data2.Current.WindKph

	// Format differences with sign
	tempDiffStr := formatDifference(tempDiff, "°C")
	feelsLikeDiffStr := formatDifference(feelsLikeDiff, "°C")
	humidityDiffStr := formatDifference(float64(humidityDiff), "%")
	windDiffStr := formatDifference(windDiff, " km/h")

	// Build table rows
	table.Append([]string{"Temperature",
		fmt.Sprintf("%.1f°C", data1.Current.TempC),
		fmt.Sprintf("%.1f°C", data2.Current.TempC),
		tempDiffStr})

	table.Append([]string{"Feels Like",
		fmt.Sprintf("%.1f°C", data1.Current.FeelsLikeC),
		fmt.Sprintf("%.1f°C", data2.Current.FeelsLikeC),
		feelsLikeDiffStr})

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

	return table
}

// formatDifference formats a numeric difference with a sign and unit
func formatDifference(diff float64, unit string) string {
	if diff == 0 {
		return "0" + unit
	}

	if diff > 0 {
		return fmt.Sprintf("+%.1f%s", diff, unit)
	}

	return fmt.Sprintf("%.1f%s", diff, unit)
}

// displayComparisonAnalysis provides textual analysis of the comparison
func displayComparisonAnalysis(data1, data2 *model.WeatherData) {
	analysisColor := color.New(color.FgHiCyan)

	// Temperature comparison
	tempDiff := data1.Current.TempC - data2.Current.TempC
	if tempDiff > 3 {
		analysisColor.Printf("📊 %s is %.1f°C warmer than %s\n",
			data1.Location.Name, tempDiff, data2.Location.Name)
	} else if tempDiff < -3 {
		analysisColor.Printf("📊 %s is %.1f°C colder than %s\n",
			data1.Location.Name, -tempDiff, data2.Location.Name)
	}

	// Humidity comparison
	humidityDiff := data1.Current.Humidity - data2.Current.Humidity
	if humidityDiff > 15 {
		analysisColor.Printf("💧 %s is more humid than %s\n",
			data1.Location.Name, data2.Location.Name)
	} else if humidityDiff < -15 {
		analysisColor.Printf("💧 %s is drier than %s\n",
			data1.Location.Name, data2.Location.Name)
	}

	// Wind comparison
	windDiff := data1.Current.WindKph - data2.Current.WindKph
	if windDiff > 10 {
		analysisColor.Printf("🌬️ %s is windier than %s\n",
			data1.Location.Name, data2.Location.Name)
	} else if windDiff < -10 {
		analysisColor.Printf("🌬️ %s is calmer than %s\n",
			data1.Location.Name, data2.Location.Name)
	}

	fmt.Println()
}
