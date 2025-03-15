// ui/chart.go
package ui

import (
	"fmt"
	"time"

	"github.com/biferdou/illapaca/model"
	"github.com/fatih/color"
)

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

// DisplayTemperatureChart renders a simple temperature chart
func DisplayTemperatureChart(data *model.WeatherData) {
	chartTitle := color.New(color.FgHiGreen, color.Bold)
	chartTitle.Println("Temperature Trend (24 hours)")
	fmt.Println()

	// Get the first day's hourly forecast
	hours := data.Forecast.ForecastDay[0].Hour

	// Process data for every 3 hours (8 points total)
	var temps []float64
	var times []string

	for i := 0; i < len(hours); i += 3 {
		if len(temps) >= 8 {
			break
		}
		temps = append(temps, hours[i].TempC)

		t, _ := time.Parse("2006-01-02 15:04", hours[i].Time)
		times = append(times, fmt.Sprintf("%02d:00", t.Hour()))
	}

	// Find min/max for scaling
	var min, max float64 = temps[0], temps[0]
	for _, t := range temps {
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}

	// Add some padding
	min = min - 1
	max = max + 1

	// Normalize temperatures to a 0-10 scale for display
	normalizedTemps := make([]int, len(temps))
	for i, t := range temps {
		// Scale to 0-10 range
		normalizedTemps[i] = int((t - min) * 10 / (max - min))
	}

	// Print the chart
	fmt.Println("     ┌────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Print("     │")

	// Top row labels
	for i, time := range times {
		if i == 0 {
			fmt.Printf("  %s", time)
		} else {
			fmt.Printf("      %s", time)
		}
	}
	fmt.Println("  │")

	// Print chart lines
	for row := 10; row >= 0; row-- {
		fmt.Print("     │")

		// Print data points
		for i, level := range normalizedTemps {
			if level == row {
				tempColor := getTemperatureColor(temps[i], min, max)
				tempColor.Print("    •      ")
			} else if level > row {
				tempColor := getTemperatureColor(temps[i], min, max)
				tempColor.Print("    │      ")
			} else {
				fmt.Print("           ")
			}
		}

		// Add temperature scale on the right side
		if row == 10 {
			fmt.Printf("│ %.1f°C", max)
		} else if row == 0 {
			fmt.Printf("│ %.1f°C", min)
		} else if row == 5 {
			midTemp := (max + min) / 2
			fmt.Printf("│ %.1f°C", midTemp)
		} else {
			fmt.Print("│")
		}

		fmt.Println()
	}

	fmt.Println("     └────────────────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
}

// DisplayPrecipitationChart renders a simple precipitation chance chart
func DisplayPrecipitationChart(day model.ForecastDay) {
	chartTitle := color.New(color.FgHiBlue, color.Bold)
	chartTitle.Println("Precipitation Chance (24 hours)")
	fmt.Println()

	// Collect data for every 3 hours
	var chances []int
	var times []string

	for i := 0; i < len(day.Hour); i += 3 {
		if len(chances) >= 8 {
			break
		}
		chances = append(chances, day.Hour[i].ChanceOfRain)

		t, _ := time.Parse("2006-01-02 15:04", day.Hour[i].Time)
		times = append(times, fmt.Sprintf("%02d:00", t.Hour()))
	}

	// Print the chart header
	fmt.Println("     ┌────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Print("     │")

	// Top row labels
	for i, time := range times {
		if i == 0 {
			fmt.Printf("  %s", time)
		} else {
			fmt.Printf("      %s", time)
		}
	}
	fmt.Println("  │")

	// Print chart
	rainLevels := []int{100, 80, 60, 40, 20, 0}
	for i, level := range rainLevels {
		fmt.Print("     │")

		// Print precipitation bars
		for _, chance := range chances {
			if chance >= level {
				// Choose symbol and color based on chance
				var symbol string
				var rainColor *color.Color

				if chance >= 80 {
					symbol = "█"
					rainColor = color.New(color.FgHiBlue, color.Bold)
				} else if chance >= 60 {
					symbol = "▓"
					rainColor = color.New(color.FgBlue)
				} else if chance >= 40 {
					symbol = "▒"
					rainColor = color.New(color.FgHiCyan)
				} else if chance >= 20 {
					symbol = "░"
					rainColor = color.New(color.FgCyan)
				} else {
					symbol = "·"
					rainColor = color.New(color.FgHiWhite)
				}

				// Print bar element with spacing
				fmt.Print("    ")
				rainColor.Print(symbol)
				fmt.Print("      ")
			} else {
				fmt.Print("           ")
			}
		}

		// Add precipitation scale on the right side
		if i == 0 {
			fmt.Print("│ 100%")
		} else if i == len(rainLevels)-1 {
			fmt.Print("│ 0%")
		} else if i == len(rainLevels)/2 {
			fmt.Print("│ 50%")
		} else {
			fmt.Print("│")
		}

		fmt.Println()
	}

	fmt.Println("     └────────────────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
}
