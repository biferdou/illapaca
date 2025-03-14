package cmd

import (
	"fmt"
	"os"

	"github.com/biferdou/illapaca/api"
	"github.com/biferdou/illapaca/ui"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current [location]",
	Short: "Show current weather conditions",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := getLocation(args)
		if location == "" {
			fmt.Println("Error: location not specified and no default location set")
			os.Exit(1)
		}

		data, err := api.FetchWeather(location, 1)
		if err != nil {
			fmt.Printf("Error fetching weather: %v\n", err)
			os.Exit(1)
		}

		ui.DisplayCurrentWeather(data)
	},
}
