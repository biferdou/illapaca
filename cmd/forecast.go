package cmd

import (
	"fmt"
	"os"

	"github.com/biferdou/illapaca/api"
	"github.com/biferdou/illapaca/ui"
	"github.com/spf13/cobra"
)

var forecastCmd = &cobra.Command{
	Use:   "forecast [location]",
	Short: "Show weather forecast",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := getLocation(args)
		if location == "" {
			fmt.Println("Error: location not specified and no default location set")
			os.Exit(1)
		}

		days, _ := cmd.Flags().GetInt("days")

		data, err := api.FetchWeather(location, days)
		if err != nil {
			fmt.Printf("Error fetching weather: %v\n", err)
			os.Exit(1)
		}

		ui.DisplayCurrentWeather(data)
		ui.DisplayForecast(data)
	},
}

func init() {
	forecastCmd.Flags().IntP("days", "d", 5, "Number of days for forecast")
}
