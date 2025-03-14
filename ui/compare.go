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
	comparisonTitle := color.New(color.FgHiBlue, color.Bold)
	comparisonTitle.Printf("Location Comparison: %s vs %s\n\n",
		data1.Location.Name, data2.Location.Name)

	// Display comparison table
	table := createComparisonTable(data1, data2)
	table.Render()
	fmt.Println()

	// Add some analysis
	displayComparisonAnalysis(data1, data2)
}

// createComparisonTable builds the comparison table for two locations
func createComparisonTable(data1, data2 *model.WeatherData) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", data1.Location.Name, data2.Location.Name, "Difference"})

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
