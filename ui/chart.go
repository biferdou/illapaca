package ui

import (
	"fmt"
	"time"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

// DisplayTemperatureChart renders an enhanced ASCII chart of temperature trends
func DisplayTemperatureChart(data *model.WeatherData) {
	chartBox := color.New(color.FgHiWhite)
	chartTitle := color.New(color.FgHiGreen, color.Bold)

	chartBox.Println("┌─────────────────────────────────────────┐")
	chartBox.Print("│ ")
	chartTitle.Println("Temperature Trend (24 hours)")
	chartBox.Println("└─────────────────────────────────────────┘")

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

	// Add a small buffer to min/max for better visualization
	maxTemp += 1
	minTemp -= 1

	chartHeight := 10
	scale := float64(chartHeight) / (maxTemp - minTemp)

	// Print temperature scale with improved styling
	maxTempColor := color.New(color.FgHiRed)
	minTempColor := color.New(color.FgHiCyan)
	axisColor := color.New(color.FgHiWhite)

	axisColor.Print("  ")
	maxTempColor.Printf("%.1f°C ", maxTemp)
	axisColor.Println("┌───────────────────────────────┐")

	// Print chart with improved styling and color gradients
	tempPoints := make([]float64, 0, len(hours)/3+1)
	timeLabels := make([]string, 0, len(hours)/3+1)

	for _, hour := range hours {
		if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
			tempPoints = append(tempPoints, hour.TempC)
			t, _ := time.Parse("2006-01-02 15:04", hour.Time)
			timeLabels = append(timeLabels, fmt.Sprintf("%02d", t.Hour()))
		}
	}

	// Print chart rows
	for i := chartHeight; i >= 0; i-- {
		threshold := minTemp + float64(i)/scale

		// Special case for middle line (add a label)
		if i == chartHeight/2 {
			midTemp := (maxTemp + minTemp) / 2
			axisColor.Print("  ")
			color.New(color.FgHiYellow).Printf("%.1f°C ", midTemp)
			axisColor.Print("│")
		} else {
			axisColor.Print("        │")
		}

		// Print data points with connecting lines
		for j, temp := range tempPoints {
			if temp >= threshold {
				// Color gradient based on temperature
				pointColor := getTemperatureColor(temp, minTemp, maxTemp)
				pointColor.Print("●")
			} else if j > 0 && tempPoints[j-1] >= threshold && temp < threshold {
				// Connecting line for downward trend
				fmt.Print("╲")
			} else if j > 0 && tempPoints[j-1] < threshold && temp >= threshold {
				// Connecting line for upward trend
				fmt.Print("╱")
			} else {
				fmt.Print(" ")
			}
		}

		axisColor.Println("│")
	}

	// Print temperature scale bottom value
	axisColor.Print("  ")
	minTempColor.Printf("%.1f°C ", minTemp)
	axisColor.Print("└")

	// Print time axis
	for range tempPoints {
		axisColor.Print("─")
	}
	axisColor.Println("┘")

	// Print hour labels
	axisColor.Print("         ")
	for _, label := range timeLabels {
		hourColor := color.New(color.FgHiBlue)
		hourColor.Printf("%s ", label)
	}
	fmt.Println()
	fmt.Println()
}

// getTemperatureColor returns a color based on temperature
func getTemperatureColor(temp, minTemp, maxTemp float64) *color.Color {
	// Calculate where this temperature falls in the range (0.0 to 1.0)
	ratio := (temp - minTemp) / (maxTemp - minTemp)

	// Create a color gradient from blue (cold) to red (hot)
	if ratio < 0.2 {
		return color.New(color.FgHiBlue)
	} else if ratio < 0.4 {
		return color.New(color.FgHiCyan)
	} else if ratio < 0.6 {
		return color.New(color.FgHiGreen)
	} else if ratio < 0.8 {
		return color.New(color.FgHiYellow)
	} else {
		return color.New(color.FgHiRed)
	}
}

// DisplayPrecipitationChart renders an ASCII chart of precipitation chances
func DisplayPrecipitationChart(day model.ForecastDay) {
	chartTitle := color.New(color.FgHiBlue, color.Bold)
	chartTitle.Println("Precipitation Chance (24 hours):")

	// Chart height and scale
	chartHeight := 7
	maxValue := 100.0 // Max percentage

	// Print chart
	for i := chartHeight; i >= 0; i-- {
		threshold := maxValue * float64(i) / float64(chartHeight)

		// Print y-axis labels on the left side
		if i == chartHeight {
			fmt.Print("100% │")
		} else if i == 0 {
			fmt.Print("  0% │")
		} else if i == chartHeight/2 {
			fmt.Print(" 50% │")
		} else {
			fmt.Print("     │")
		}

		// Print data points
		for _, hour := range day.Hour {
			if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
				chanceOfRain := float64(hour.ChanceOfRain)
				if chanceOfRain >= threshold {
					// Use different symbols based on precipitation chance
					if chanceOfRain > 70 {
						fmt.Print("█")
					} else if chanceOfRain > 40 {
						fmt.Print("▓")
					} else if chanceOfRain > 10 {
						fmt.Print("▒")
					} else {
						fmt.Print("░")
					}
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}

	// Print time axis
	fmt.Print("     └")
	for _, hour := range day.Hour {
		if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
			fmt.Print("─")
		}
	}
	fmt.Println()

	// Print hours
	fmt.Print("       ")
	for _, hour := range day.Hour {
		if hour.TimeEpoch%(3600*3) == 0 { // Every 3 hours
			t, _ := time.Parse("2006-01-02 15:04", hour.Time)
			fmt.Printf("%02d ", t.Hour())
		}
	}
	fmt.Println()
	fmt.Println()
}
