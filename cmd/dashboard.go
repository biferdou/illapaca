package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/biferdou/illapaca/api"
	"github.com/biferdou/illapaca/ui"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard [location]",
	Short: "Show complete weather dashboard",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := getLocation(args)
		if location == "" {
			fmt.Println("Error: location not specified and no default location set")
			os.Exit(1)
		}

		data, err := api.FetchWeather(location, 5)
		if err != nil {
			fmt.Printf("Error fetching weather: %v\n", err)
			os.Exit(1)
		}

		ui.DisplayDashboard(data)

		// Listen for Ctrl+C to exit dashboard
		fmt.Println("Press Ctrl+C to exit dashboard")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
	},
}
