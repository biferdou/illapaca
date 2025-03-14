// cmd/compare.go
package cmd

import (
	"fmt"
	"os"

	"github.com/biferdou/illapaca/api"
	"github.com/biferdou/illapaca/ui"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare [location1] [location2]",
	Short: "Compare weather between two different locations",
	Long:  `Compare current weather conditions between two different locations side by side.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		location1 := args[0]
		location2 := args[1]

		// Get weather for first location
		data1, err := api.FetchWeather(location1, 1)
		if err != nil {
			fmt.Printf("Error fetching weather for %s: %v\n", location1, err)
			os.Exit(1)
		}

		// Get weather for second location
		data2, err := api.FetchWeather(location2, 1)
		if err != nil {
			fmt.Printf("Error fetching weather for %s: %v\n", location2, err)
			os.Exit(1)
		}

		// Display comparison
		ui.DisplayLocationComparison(data1, data2)
	},
}

func init() {
	// No flags needed for this command
}
