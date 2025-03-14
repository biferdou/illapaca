package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/biferdou/illapaca/api"
	"github.com/biferdou/illapaca/ui"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare [location]",
	Short: "Compare current weather with historical data",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := getLocation(args)
		if location == "" {
			fmt.Println("Error: location not specified and no default location set")
			os.Exit(1)
		}

		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			// Default to same day last year
			now := time.Now()
			lastYear := now.Year() - 1
			date = fmt.Sprintf("%d-%02d-%02d", lastYear, now.Month(), now.Day())
		}

		// Get current weather
		current, err := api.FetchWeather(location, 1)
		if err != nil {
			fmt.Printf("Error fetching current weather: %v\n", err)
			os.Exit(1)
		}

		// Get historical weather
		historical, err := api.FetchHistoricalWeather(location, date)
		if err != nil {
			fmt.Printf("Error fetching historical weather: %v\n", err)
			os.Exit(1)
		}

		ui.DisplayCurrentWeather(current)
		ui.DisplayHistoricalComparison(current, historical)
	},
}

func init() {
	compareCmd.Flags().StringP("date", "d", "", "Date for historical comparison (YYYY-MM-DD)")
}
