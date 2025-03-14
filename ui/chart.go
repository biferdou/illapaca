package ui

import (
	"fmt"
	"time"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

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
